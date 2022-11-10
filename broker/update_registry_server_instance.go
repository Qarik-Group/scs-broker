package broker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/lager"
	brokerapi "github.com/pivotal-cf/brokerapi/domain"
	"github.com/starkandwayne/scs-broker/broker/utilities"
	scsccparser "github.com/starkandwayne/spring-cloud-services-cli-config-parser"
)

func (broker *SCSBroker) updateRegistryServerInstance(cxt context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	spec := brokerapi.UpdateServiceSpec{}

	appName := utilities.MakeAppName(details.ServiceID, instanceID)
	spaceGUID := broker.Config.InstanceSpaceGUID

	broker.Logger.Info("update-service-instance", lager.Data{"plan-id": details.PlanID, "service-id": details.ServiceID})
	envsetup := scsccparser.EnvironmentSetup{}
	cfClient, err := broker.GetClient()
	if err != nil {
		return spec, errors.New("Couldn't start session: " + err.Error())
	}

	community, err := broker.GetCommunity()
	if err != nil {
		return spec, err
	}

	rc := utilities.NewRegistryConfig()
	rp, err := utilities.ExtractRegistryParams(string(details.RawParameters))
	if err != nil {
		return spec, err
	}

	count, err := rp.Count()
	if err != nil {
		return spec, err
	}

	info, _, _, err := cfClient.GetInfo()
	if err != nil {
		return spec, err
	}

	app, _, err := cfClient.GetApplicationByNameAndSpace(appName, spaceGUID)
	if err != nil {
		return spec, errors.New("Couldn't find app session: " + err.Error())
	}

	mapparams, err := envsetup.ParseEnvironmentFromRaw(details.RawParameters)
	if err != nil {
		return spec, err
	}

	broker.Logger.Info("Updating Environment")
	err = broker.UpdateAppEnvironment(cfClient, &app, &info, details.ServiceID, instanceID, string(details.RawParameters), mapparams)
	if err != nil {
		return spec, err
	}

	broker.Logger.Info("Updating application")

	_, _, err = cfClient.UpdateApplication(utilities.SafeApp(app))
	if err != nil {
		broker.Logger.Info("UpdateApplication(app) failed")
		return spec, err
	}

	broker.Logger.Info("handling node count")
	// handle the node count
	if count > 1 {
		rc.Clustered()
	} else {
		rc.Standalone()
	}

	// since this is an update, we need to scale, but only if the desired proc
	// count has changed
	procs, err := getApplicationProcessesByType(cfClient, broker.Logger, app.GUID, "web")
	if err != nil {
		return spec, err
	}

	procCount := 0
	for _, proc := range procs {
		if proc.Instances.IsSet {
			procCount += proc.Instances.Value
		}
	}

	broker.Logger.Info(fmt.Sprintf("I received %d procs from the API", procCount))

	if count != procCount {
		broker.Logger.Info(fmt.Sprintf("Scaling to %d procs", count))
		err = broker.scaleRegistryServer(cfClient, &app, count)
		if err != nil {
			return spec, err
		}
	}

	if count > 1 {
		stats, err := getProcessStatsByAppAndType(cfClient, community, broker.Logger, app.GUID, "web")
		if err != nil {
			return spec, err
		}

		for _, stat := range stats {
			rc.AddPeer(stat.Index, "http", stat.Host, stat.InstancePorts[0].External)
		}
	}

	broker.Logger.Info("Updating Environment")
	err = broker.UpdateRegistryEnvironment(cfClient, &app, &info, details.ServiceID, instanceID, rc, mapparams)

	if err != nil {
		return spec, err
	}

	/*
		app, _, err = cfClient.UpdateApplicationRestart(app.GUID)
		if err != nil {
			return spec, err
		}*/

	domains, _, err := cfClient.GetDomains(
		ccv3.Query{Key: ccv3.NameFilter, Values: []string{broker.Config.InstanceDomain}},
	)
	if err != nil {
		return spec, err
	}

	route, _, err := cfClient.CreateRoute(ccv3.Route{
		SpaceGUID:  spaceGUID,
		DomainGUID: domains[0].GUID,
		Host:       appName,
	})
	if err != nil {
		return spec, err
	}

	_, err = cfClient.MapRoute(route.GUID, app.GUID)

	peers, err := json.Marshal(rc.Peers)
	if err != nil {
		return spec, err
	}
	for _, peer := range rc.Peers {
		req, err := http.NewRequest(http.MethodPost, "https://"+route.URL+"/cf-config/peers", bytes.NewBuffer(peers))
		if err != nil {
			fmt.Printf("client: could not create request: %s\n", err)

		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Cf-App-Instance", app.GUID+":"+string(peer.Index))

		client := http.Client{
			Timeout: 30 * time.Second,
		}

		res, err := client.Do(req)
		if err != nil {
			fmt.Printf("client: error making http request: %s\n", err)
		}
		broker.Logger.Info(res.Status)
	}

	return spec, nil
}

func getApplicationProcessesByType(client *ccv3.Client, logger lager.Logger, appGUID string, procType string) ([]ccv3.Process, error) {
	filtered := make([]ccv3.Process, 0)

	candidates, _, err := client.GetApplicationProcesses(appGUID)
	if err != nil {
		return filtered, err
	}

	logger.Info(fmt.Sprintf("getApplicationProcessesByType got %d total procs", len(candidates)))

	for _, prospect := range candidates {

		if prospect.Type == procType {
			filtered = append(filtered, prospect)
		}
	}

	return filtered, nil
}
