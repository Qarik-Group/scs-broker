package configserver

import (
	"errors"
	"fmt"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3/constant"
	"github.com/starkandwayne/scs-broker/broker/utilities"
	"github.com/starkandwayne/scs-broker/client"
	"github.com/starkandwayne/scs-broker/config"
	"github.com/starkandwayne/scs-broker/fs"
	"github.com/starkandwayne/scs-broker/logger"
	"github.com/starkandwayne/scs-broker/poll"
)

func (broker *Broker) createInstance(serviceId string, instanceId string, jsonparams string, params map[string]string) (string, error) {

	service, err := config.GetServiceByServiceID(serviceId)
	if err != nil {
		return "", err
	}
	cfClient, err := client.GetClient()
	if err != nil {
		return "", errors.New("Couldn't start session: " + err.Error())
	}
	appName := utilities.MakeAppName(serviceId, instanceId)
	spaceGUID := config.Parsed.InstanceSpaceGUID

	logger.Info("Creating Application")
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

	logger.Info("Updating Environment")
	err = broker.UpdateAppEnvironment(cfClient, &app, &info, serviceId, instanceId, jsonparams, params)

	if err != nil {
		return "", err
	}

	logger.Info("Creating Package")
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

	logger.Info("Uploading Package")

	artifact, fi, err := fs.ArtifactStat(service.ServiceDownloadURI)
	if err != nil {
		return "", err
	}

	logger.Info(fmt.Sprintf("Uploadinlsg: %s from %s size(%d)", fi.Name(), artifact, fi.Size()))

	upkg, uwarnings, err := cfClient.UploadPackage(pkg, artifact)
	logger.ShowWarnings(uwarnings, upkg)
	if err != nil {
		return "", err
	}

	logger.Info("Polling Package")
	pkg, pwarnings, err := poll.Package(pkg)
	logger.ShowWarnings(pwarnings, pkg)
	if err != nil {

		return "", err
	}

	logger.Info("Creating Build")
	build, cwarnings, err := cfClient.CreateBuild(ccv3.Build{PackageGUID: pkg.GUID})
	logger.ShowWarnings(cwarnings, build)
	if err != nil {
		return "", err
	}

	logger.Info("polling build")
	droplet, pbwarnings, err := poll.Build(build.GUID, appName)
	logger.ShowWarnings(pbwarnings, droplet)
	if err != nil {
		return "", err
	}

	logger.Info("set application droplet")
	_, _, err = cfClient.SetApplicationDroplet(app.GUID, droplet.GUID)
	if err != nil {
		return "", err
	}
	domains, _, err := cfClient.GetDomains(
		ccv3.Query{Key: ccv3.NameFilter, Values: []string{config.Parsed.InstanceDomain}},
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
	app, _, err = cfClient.UpdateApplicationRestart(app.GUID)
	if err != nil {
		return "", err
	}

	logger.Info(route.URL)

	return route.URL, nil
}
