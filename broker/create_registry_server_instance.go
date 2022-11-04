package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3/constant"
	"github.com/starkandwayne/scs-broker/broker/utilities"
)

func (broker *SCSBroker) createRegistryServerInstance(serviceId string, instanceId string, jsonparams string, params map[string]string) (string, error) {
	service, err := broker.GetServiceByServiceID(serviceId)
	if err != nil {
		return "", err
	}

	rc := &registryConfig{}
	broker.Logger.Info("jsonparams == " + jsonparams)
	rp, err := utilities.ExtractRegistryParams(jsonparams)
	if err != nil {
		return "", err
	}

	count, err := rp.Count()
	if err != nil {
		return "", err
	}

	cfClient, err := broker.GetClient()
	if err != nil {
		return "", errors.New("Couldn't start session: " + err.Error())
	}
	appName := utilities.MakeAppName(serviceId, instanceId)
	spaceGUID := broker.Config.InstanceSpaceGUID

	broker.Logger.Info("Creating Application")
	app, _, err := cfClient.CreateApplication(
		ccv3.Application{
			Name:          appName,
			LifecycleType: constant.AppLifecycleTypeBuildpack,
			State:         constant.ApplicationStopped,
			Relationships: ccv3.Relationships{
				constant.RelationshipTypeSpace: ccv3.Relationship{GUID: spaceGUID},
			},
		},
	)
	if err != nil {
		return "", err
	}

	info, _, _, err := cfClient.GetInfo()
	if err != nil {
		return "", err
	}

	broker.Logger.Info("Updating Environment")
	err = broker.UpdateAppEnvironment(cfClient, &app, &info, serviceId, instanceId, jsonparams, params)

	if err != nil {
		return "", err
	}

	broker.Logger.Info("Creating Package")
	pkg, _, err := cfClient.CreatePackage(
		ccv3.Package{
			Type: constant.PackageTypeBits,
			Relationships: ccv3.Relationships{
				constant.RelationshipTypeApplication: ccv3.Relationship{GUID: app.GUID},
			},
		})
	if err != nil {
		return "", err
	}

	broker.Logger.Info("Uploading Package")

	jarname := path.Base(service.ServiceDownloadURI)
	artifact := broker.Config.ArtifactsDir + "/" + jarname

	fi, err := os.Stat(artifact)
	if err != nil {
		return "", err
	}

	broker.Logger.Info(fmt.Sprintf("Uploading: %s from %s size(%d)", fi.Name(), artifact, fi.Size()))

	upkg, uwarnings, err := cfClient.UploadPackage(pkg, artifact)
	broker.showWarnings(uwarnings, upkg)
	if err != nil {
		return "", err
	}

	broker.Logger.Info("Polling Package")
	pkg, pwarnings, err := broker.pollPackage(pkg)
	broker.showWarnings(pwarnings, pkg)
	if err != nil {

		return "", err
	}

	broker.Logger.Info("Creating Build")
	build, cwarnings, err := cfClient.CreateBuild(ccv3.Build{PackageGUID: pkg.GUID})
	broker.showWarnings(cwarnings, build)
	if err != nil {
		return "", err
	}

	broker.Logger.Info("polling build")
	droplet, pbwarnings, err := broker.pollBuild(build.GUID, appName)
	broker.showWarnings(pbwarnings, droplet)
	if err != nil {
		return "", err
	}

	broker.Logger.Info("set application droplet")
	_, _, err = cfClient.SetApplicationDroplet(app.GUID, droplet.GUID)
	if err != nil {
		return "", err
	}
	domains, _, err := cfClient.GetDomains(
		ccv3.Query{Key: ccv3.NameFilter, Values: []string{broker.Config.InstanceDomain}},
	)
	if err != nil {
		return "", err
	}

	if len(domains) == 0 {
		return "", errors.New("no domains found for this instance")
	}

	route, _, err := cfClient.CreateRoute(ccv3.Route{
		SpaceGUID:  spaceGUID,
		DomainGUID: domains[0].GUID,
		Host:       appName,
	})
	if err != nil {
		return "", err
	}
	_, err = cfClient.MapRoute(route.GUID, app.GUID)
	if err != nil {
		return "", err
	}

	broker.Logger.Info("handle node count")
	// handle the node count
	if count > 1 {
		rc.Clustered()
		rc.AddPeer(fmt.Sprintf("%s/eureka", route.URL))
		broker.Logger.Info(fmt.Sprintf("scaling to %d", count))
		err = broker.scaleRegistryServer(cfClient, &app, count, rc)
		if err != nil {
			return "", err
		}
	} else {
		rc.Standalone()
	}

	broker.Logger.Info("Updating Environment")
	err = broker.UpdateAppEnvironment(cfClient, &app, &info, serviceId, instanceId, rc.String(), params)

	if err != nil {
		return "", err
	}

	app, _, err = cfClient.UpdateApplicationRestart(app.GUID)
	if err != nil {
		return "", err
	}

	broker.Logger.Info(route.URL)

	return route.URL, nil
}

type registryConfig struct {
	Mode  string
	Peers []string
}

func (rc *registryConfig) AddPeer(peer string) {
	rc.Peers = append(rc.Peers, peer)
}

func (rc *registryConfig) Standalone() {
	rc.Mode = "standalone"
}

func (rc *registryConfig) Clustered() {
	rc.Mode = "clustered"
}

func (rc *registryConfig) String() string {
	return string(rc.Bytes())
}

func (rc *registryConfig) Bytes() []byte {
	client := make(map[string]interface{})

	if rc.Mode == "standalone" {
		client["registerWithEureka"] = false
		client["fetchRegistry"] = false
	}

	if len(rc.Peers) > 0 {
		serviceUrl := make(map[string]interface{})
		defaultZone := strings.Join(rc.Peers, ",")
		serviceUrl["defaultZone"] = defaultZone
		client["serviceUrl"] = serviceUrl
	}

	eureka := make(map[string]interface{})
	eureka["client"] = client

	data := make(map[string]interface{})
	data["eureka"] = eureka

	output, err := json.Marshal(data)
	if err != nil {
		return []byte("{}")
	}

	return output

}
