package broker

import (
	"errors"
	"fmt"
	"os"
	"path"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3/constant"
	"github.com/starkandwayne/scs-broker/broker/utilities"
)

func (broker *SCSBroker) createRegistryServerInstance(serviceId string, instanceId string, jsonparams string, params map[string]string) (string, error) {
	service, err := broker.GetServiceByServiceID(serviceId)
	if err != nil {
		return "", err
	}

	rc := utilities.NewRegistryConfig()
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
		broker.Logger.Info(fmt.Sprintf("scaling to %d", count))
		err = broker.scaleRegistryServer(cfClient, &app, count)
		if err != nil {
			return "", err
		}

		community, err := broker.GetCommunity()
		if err != nil {
			return "", err
		}

		stats, err := getProcessStatsByAppAndType(cfClient, community, broker.Logger, app.GUID, "web")
		if err != nil {
			return "", nil
		}

		for _, stat := range stats {
			rc.AddPeer(stat.Index, "http", stat.Host, stat.InstancePorts[0].External)
		}
	} else {
		rc.Standalone()
	}

	broker.Logger.Info("Updating Environment")
	err = broker.UpdateRegistryEnvironment(cfClient, &app, &info, serviceId, instanceId, rc, params)

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
