package broker

import (
	"context"
	"errors"

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
	_, _, err = cfClient.UpdateApplication(app)
	if err != nil {
		return spec, err
	}

	broker.Logger.Info("handling node count")
	// handle the node count
	rc := &registryConfig{}
	rp, err := utilities.ExtractRegistryParams(string(details.RawParameters))
	if err != nil {
		return spec, err
	}

	if count, found := rp["count"]; found {
		if c, ok := count.(int); ok {
			if c > 1 {
				rc.Clustered()
				err = broker.scaleRegistryServer(cfClient, &app, c, rc)
				if err != nil {
					return spec, err
				}
			} else {
				rc.Standalone()
			}
		} else {
			rc.Standalone()
		}
	}

	broker.Logger.Info("Updating Environment")
	err = broker.UpdateAppEnvironment(cfClient, &app, &info, details.ServiceID, instanceID, rc.String(), mapparams)

	if err != nil {
		return spec, err
	}

	_, _, err = cfClient.UpdateApplicationRestart(app.GUID)
	if err != nil {
		return spec, err
	}

	return spec, nil
}
