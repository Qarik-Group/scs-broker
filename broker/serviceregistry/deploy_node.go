package serviceregistry

import (
	"errors"
	"fmt"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3/constant"
	"github.com/starkandwayne/scs-broker/broker/serviceregistry/result"
	"github.com/starkandwayne/scs-broker/client"
	"github.com/starkandwayne/scs-broker/config"
	"github.com/starkandwayne/scs-broker/fs"
	"github.com/starkandwayne/scs-broker/logger"
	"github.com/starkandwayne/scs-broker/poll"
)

func (broker *Broker) deployNode(serviceId string, appName string, pipeline chan<- *result.PushApp) {
	service, err := config.GetServiceByServiceID(serviceId)
	if err != nil {
		pipeline <- result.NewPushApp().WithError(err)
		return
	}

	cfClient, err := client.GetClient()
	if err != nil {
		pipeline <- result.NewPushApp().WithError(err)
		return
	}

	spaceGUID := config.Parsed.InstanceSpaceGUID

	broker.logDeployNodeInfo("Creating Application", appName, spaceGUID)
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
		pipeline <- result.NewPushApp().WithError(err)
		return
	}

	broker.logDeployNodeInfo("Updating Environment", appName, spaceGUID)

	broker.logDeployNodeInfo("Creating Package", appName, spaceGUID)
	pkg, _, err := cfClient.CreatePackage(
		ccv3.Package{
			Type: constant.PackageTypeBits,
			Relationships: ccv3.Relationships{
				constant.RelationshipTypeApplication: ccv3.Relationship{GUID: app.GUID},
			},
		})
	if err != nil {
		pipeline <- result.NewPushApp().WithError(err)
		return
	}

	broker.logDeployNodeInfo("Uploading Package", appName, spaceGUID)

	artifact, fi, err := fs.ArtifactStat(service.ServiceDownloadURI)
	if err != nil {
		pipeline <- result.NewPushApp().WithError(err)
		return
	}

	broker.logDeployNodeInfo(
		fmt.Sprintf("Uploading: %s from %s size(%d)", fi.Name(), artifact, fi.Size()),
		appName,
		spaceGUID,
	)

	upkg, uwarnings, err := cfClient.UploadPackage(pkg, artifact)
	logger.ShowWarnings(uwarnings, upkg)
	if err != nil {
		pipeline <- result.NewPushApp().WithError(err)
		return
	}

	broker.logDeployNodeInfo("Rolling Package", appName, spaceGUID)
	pkg, pwarnings, err := poll.Package(pkg)
	logger.ShowWarnings(pwarnings, pkg)
	if err != nil {
		pipeline <- result.NewPushApp().WithError(err)
		return
	}

	broker.logDeployNodeInfo("Creating Build", appName, spaceGUID)
	build, cwarnings, err := cfClient.CreateBuild(ccv3.Build{PackageGUID: pkg.GUID})
	logger.ShowWarnings(cwarnings, build)
	if err != nil {
		pipeline <- result.NewPushApp().WithError(err)
		return
	}

	broker.logDeployNodeInfo("polling build", appName, spaceGUID)
	droplet, pbwarnings, err := poll.Build(build.GUID, appName)
	logger.ShowWarnings(pbwarnings, droplet)
	if err != nil {
		pipeline <- result.NewPushApp().WithError(err)
		return
	}

	broker.logDeployNodeInfo("set application droplet", appName, spaceGUID)
	_, _, err = cfClient.SetApplicationDroplet(app.GUID, droplet.GUID)
	if err != nil {
		pipeline <- result.NewPushApp().WithError(err)
		return
	}

	domains, _, err := cfClient.GetDomains(
		ccv3.Query{Key: ccv3.NameFilter, Values: []string{config.Parsed.InstanceDomain}},
	)
	if err != nil {
		pipeline <- result.NewPushApp().WithError(err)
		return
	}

	if len(domains) == 0 {
		pipeline <- result.NewPushApp().WithError(errors.New("no domains found for this instance"))
		return
	}

	route, _, err := cfClient.CreateRoute(ccv3.Route{
		SpaceGUID:  spaceGUID,
		DomainGUID: domains[0].GUID,
		Host:       appName,
	})
	if err != nil {
		pipeline <- result.NewPushApp().WithError(err)
		return
	}
	_, err = cfClient.MapRoute(route.GUID, app.GUID)
	if err != nil {
		pipeline <- result.NewPushApp().WithError(err)
		return
	}

	app, _, err = cfClient.UpdateApplicationRestart(app.GUID)
	if err != nil {
		pipeline <- result.NewPushApp().WithError(err)
		return
	}

	pipeline <- result.NewPushApp().WithApp(app).WithRoute(route)
}
