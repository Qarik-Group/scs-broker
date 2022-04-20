package broker

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccerror"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3/constant"
	"code.cloudfoundry.org/cli/types"
	"code.cloudfoundry.org/cli/util/configv3"
	"code.cloudfoundry.org/lager"
	"github.com/cloudfoundry-community/go-uaa"
	brokerapi "github.com/pivotal-cf/brokerapi/domain"
	"github.com/starkandwayne/config-server-broker/config"
	scsccparser "github.com/starkandwayne/spring-cloud-services-cli-config-parser"
)

const (
	ArtifactsDir string = "artifacts"
)

type ConfigServerBroker struct {
	Config config.Config
	Logger lager.Logger
}

func (broker *ConfigServerBroker) Services(ctx context.Context) ([]brokerapi.Service, error) {
	planList := []brokerapi.ServicePlan{
		{
			ID:          broker.Config.BasicPlanId,
			Name:        broker.Config.BasicPlanName,
			Description: "This plan provides a config server deployed to cf",
			Metadata: &brokerapi.ServicePlanMetadata{
				DisplayName: "Basic",
			},
		}}

	return []brokerapi.Service{
		{
			ID:          broker.Config.ServiceID,
			Name:        broker.Config.ServiceName,
			Description: broker.Config.Description,
			Bindable:    true,
			Plans:       planList,
			Metadata: &brokerapi.ServiceMetadata{
				DisplayName:         broker.Config.DisplayName,
				LongDescription:     broker.Config.LongDescription,
				DocumentationUrl:    broker.Config.DocumentationURL,
				SupportUrl:          broker.Config.SupportURL,
				ImageUrl:            fmt.Sprintf("data:image/png;base64,%s", broker.Config.IconImage),
				ProviderDisplayName: broker.Config.ProviderDisplayName,
			},
			Tags: []string{
				"snw",
				"config-server",
			},
		},
	}, nil
}

func (broker *ConfigServerBroker) Provision(ctx context.Context, instanceID string, serviceDetails brokerapi.ProvisionDetails, asyncAllowed bool) (spec brokerapi.ProvisionedServiceSpec, err error) {
	spec = brokerapi.ProvisionedServiceSpec{}
	envsetup := scsccparser.EnvironmentSetup{}
	mapparams, err := envsetup.ParseEnvironmentFromRaw(serviceDetails.RawParameters)
	if err != nil {
		return spec, err
	}

	if serviceDetails.PlanID != broker.Config.BasicPlanId {
		return spec, errors.New("plan_id not recognized")
	}

	err = broker.createBasicInstance(instanceID, mapparams)
	if err != nil {
		return spec, err
	}

	return spec, nil
}

func (broker *ConfigServerBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	spec := brokerapi.DeprovisionServiceSpec{}
	cfClient, err := broker.getClient()
	if err != nil {
		return spec, err
	}
	appName := makeAppName(instanceID)
	app, _, err := cfClient.GetApplicationByNameAndSpace(appName, broker.Config.InstanceSpaceGUID)
	appNotFound := ccerror.ApplicationNotFoundError{Name: appName}
	if err == appNotFound {
		broker.Logger.Info("app-not-found")
		return spec, nil
	}

	if err != nil {
		return spec, err
	}
	routes, _, err := cfClient.GetApplicationRoutes(app.GUID)
	if err != nil {
		return spec, err
	}
	_, _, err = cfClient.UpdateApplicationStop(app.GUID)
	if err != nil {
		return spec, err
	}

	for route := range routes {
		_, _, err := cfClient.DeleteRoute(routes[route].GUID)
		if err != nil {
			return spec, err
		}
	}

	_, _, err = cfClient.DeleteApplication(app.GUID)
	if err != nil {
		return spec, err
	}

	return spec, nil
}

func (broker *ConfigServerBroker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails, asyncAllowed bool) (brokerapi.UnbindSpec, error) {
	unbind := brokerapi.UnbindSpec{}
	broker.Logger.Info("UnBind: GetUAAClient")
	api, err := broker.getUaaClient()
	if err != nil {
		broker.Logger.Info("UnBind: Error in GetUAAClient")
		return unbind, err
	}

	broker.Logger.Info("UnBind: makeClientIdForBinding")
	clientId := broker.makeClientIdForBinding(bindingID)

	broker.Logger.Info(fmt.Sprintf("UnBind: DeleteClient bindingID:%s clientid %s", bindingID, clientId))
	_, err = api.DeleteClient(clientId)
	if err != nil {
		broker.Logger.Error("UnBind: Error in DeleteClient - will attempt to remove anyway", err)
		return unbind, nil
	}
	broker.Logger.Info("UnBind: Return")
	return unbind, nil
}

func (broker *ConfigServerBroker) makeClientIdForBinding(bindingId string) string {
	return "config-server-binding-" + strings.Replace(bindingId, broker.Config.ServiceID+"-", "", 1)
}
func (broker *ConfigServerBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (brokerapi.Binding, error) {
	binding := brokerapi.Binding{}
	broker.Logger.Info("Bind: GetUAAClient")
	api, err := broker.getUaaClient()
	if err != nil {
		broker.Logger.Info("Bind: Error in getting client")
		return binding, err
	}
	clientId := broker.makeClientIdForBinding(bindingID)
	password := broker.genClientPassword()
	client := uaa.Client{
		ClientID:             clientId,
		AuthorizedGrantTypes: []string{"client_credentials"},
		Authorities:          []string{fmt.Sprintf("config-server.%v.read", instanceID)},
		DisplayName:          clientId,
		ClientSecret:         password,
	}

	broker.Logger.Info("Bind: got client info")
	broker.Logger.Info("Bind: Create Client")
	_, err = api.CreateClient(client)
	if err != nil {
		broker.Logger.Info("Bind: Error in CreateClient")
		return binding, err
	}

	broker.Logger.Info("Bind: GetClient")
	cfClient, err := broker.getClient()
	if err != nil {
		broker.Logger.Info("Bind: Error in GetClient")
		return binding, err
	}

	broker.Logger.Info("Bind: Get Info")
	info, _, _, err := cfClient.GetInfo()
	if err != nil {
		broker.Logger.Info("Bind: Error in Get Info")

		return binding, err
	}

	broker.Logger.Info("Bind: GetApplicationByNameAndSpace")

	app, _, err := cfClient.GetApplicationByNameAndSpace(makeAppName(instanceID), broker.Config.InstanceSpaceGUID)
	if err != nil {
		broker.Logger.Info("Bind: Error in GetApplicationByNameAndSpace")
		return binding, err
	}

	broker.Logger.Info("Bind: GetApplicationRoutes")
	routes, _, err := cfClient.GetApplicationRoutes(app.GUID)
	if err != nil {
		broker.Logger.Info("Bind: Error in GetApplicationRoutes")
		return binding, err
	}

	broker.Logger.Info("Bind: Building binding Credentials")
	binding.Credentials = map[string]string{
		"uri":              fmt.Sprintf("https://%v", routes[0].URL),
		"access_token_uri": fmt.Sprintf("%v/oauth/token", info.UAA()),
		"client_id":        clientId,
		"client_secret":    password,
	}

	broker.Logger.Info("Bind: Return")

	return binding, nil
}

// LastOperation ...
// If the broker provisions asynchronously, the Cloud Controller will poll this endpoint
// for the status of the provisioning operation.
func (broker *ConfigServerBroker) LastOperation(ctx context.Context, instanceID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, errors.New("not implemented")
}

func (broker *ConfigServerBroker) Update(cxt context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	return brokerapi.UpdateServiceSpec{}, errors.New("not implemented")
}

func (broker *ConfigServerBroker) GetBinding(ctx context.Context, instanceID, bindingID string) (brokerapi.GetBindingSpec, error) {
	return brokerapi.GetBindingSpec{}, errors.New("not implemented")
}

func (broker *ConfigServerBroker) GetInstance(ctx context.Context, instanceID string) (brokerapi.GetInstanceDetailsSpec, error) {
	return brokerapi.GetInstanceDetailsSpec{}, errors.New("not implemented")
}

func (broker *ConfigServerBroker) LastBindingOperation(ctx context.Context, instanceID, bindingID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	//create client

	return brokerapi.LastOperation{}, errors.New("not implemented")
}

func makeAppName(instanceId string) string {
	return "config-server-" + instanceId
}
func (broker *ConfigServerBroker) createBasicInstance(instanceId string, params map[string]string) error {
	cfClient, err := broker.getClient()
	if err != nil {
		return errors.New("Couldn't start session: " + err.Error())
	}
	appName := makeAppName(instanceId)
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
		return err
	}

	info, _, _, err := cfClient.GetInfo()
	if err != nil {
		return err
	}

	for key, value := range params {
		_, _, err := cfClient.UpdateApplicationEnvironmentVariables(app.GUID, ccv3.EnvironmentVariables{
			key: *types.NewFilteredString(value),
		})
		if err != nil {
			return err
		}
	}

	_, _, err = cfClient.UpdateApplicationEnvironmentVariables(app.GUID, ccv3.EnvironmentVariables{
		//"SPRING_CLOUD_CONFIG_SERVER_GIT_URI": *types.NewFilteredString(params.GitRepoUrl),
		//"JBP_CONFIG_OPEN_JDK_JRE": *types.NewFilteredString("{ jre: { version: 8.+ } }"),
		"JWK_SET_URI":         *types.NewFilteredString(fmt.Sprintf("%v/token_keys", info.UAA())),
		"SKIP_SSL_VALIDATION": *types.NewFilteredString(strconv.FormatBool(broker.Config.CfConfig.SkipSslValidation)),
		"REQUIRED_AUDIENCE":   *types.NewFilteredString(fmt.Sprintf("config-server.%v", instanceId)),
	})

	if err != nil {
		return err
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
		return err
	}

	broker.Logger.Info("Uploading Package")

	artifact := "./" + ArtifactsDir + "/spring-cloud-config-server.jar"

	fi, err := os.Stat(artifact)
	if err != nil {
		return err
	}

	broker.Logger.Info(fmt.Sprintf("Uploading: %s from %s size(%d)", fi.Name(), artifact, fi.Size()))

	upkg, uwarnings, err := cfClient.UploadPackage(pkg, artifact)
	broker.showWarnings(uwarnings, upkg)
	if err != nil {
		return err
	}

	broker.Logger.Info("Polling Package")
	pkg, pwarnings, err := broker.pollPackage(pkg)
	broker.showWarnings(pwarnings, pkg)
	if err != nil {

		return err
	}

	broker.Logger.Info("Creating Build")
	build, cwarnings, err := cfClient.CreateBuild(ccv3.Build{PackageGUID: pkg.GUID})
	broker.showWarnings(cwarnings, build)
	if err != nil {
		return err
	}

	broker.Logger.Info("polling build")
	droplet, pbwarnings, err := broker.pollBuild(build.GUID, appName)
	broker.showWarnings(pbwarnings, droplet)
	if err != nil {
		return err
	}

	broker.Logger.Info("set application droplet")
	_, _, err = cfClient.SetApplicationDroplet(app.GUID, droplet.GUID)
	if err != nil {
		return err
	}
	domains, _, err := cfClient.GetDomains(
		ccv3.Query{Key: ccv3.NameFilter, Values: []string{broker.Config.InstanceDomain}},
	)
	if err != nil {
		return err
	}

	if len(domains) == 0 {
		return errors.New("no domains found for this instance")
	}

	route, _, err := cfClient.CreateRoute(ccv3.Route{
		SpaceGUID:  spaceGUID,
		DomainGUID: domains[0].GUID,
		Host:       appName,
	})
	if err != nil {
		return err
	}
	_, err = cfClient.MapRoute(route.GUID, app.GUID)
	if err != nil {
		return err
	}
	_, _, err = cfClient.UpdateApplicationRestart(app.GUID)
	if err != nil {
		return err
	}

	return nil
}

func (broker *ConfigServerBroker) pollBuild(buildGUID string, appName string) (ccv3.Droplet, ccv3.Warnings, error) {
	var allWarnings ccv3.Warnings

	timeout := time.After(configv3.DefaultStagingTimeout)
	interval := time.NewTimer(0)

	cfClient, err := broker.getClient()
	if err != nil {
		return ccv3.Droplet{}, nil, errors.New("couldn't start session: " + err.Error())
	}

	for {
		select {
		case <-interval.C:
			build, warnings, err := cfClient.GetBuild(buildGUID)
			allWarnings = append(allWarnings, warnings...)
			if err != nil {
				return ccv3.Droplet{}, allWarnings, err
			}

			switch build.State {
			case constant.BuildFailed:
				return ccv3.Droplet{}, allWarnings, errors.New(build.Error)

			case constant.BuildStaged:
				droplet, warnings, err := cfClient.GetDroplet(build.DropletGUID)
				allWarnings = append(allWarnings, warnings...)
				if err != nil {
					return ccv3.Droplet{}, allWarnings, err
				}

				return ccv3.Droplet{
					GUID:      droplet.GUID,
					State:     droplet.State,
					CreatedAt: droplet.CreatedAt,
				}, allWarnings, nil
			}

			interval.Reset(configv3.DefaultPollingInterval)

		case <-timeout:
			return ccv3.Droplet{}, allWarnings, errors.New("staging timed out")
		}
	}
}

func (broker *ConfigServerBroker) pollPackage(pkg ccv3.Package) (ccv3.Package, ccv3.Warnings, error) {
	var allWarnings ccv3.Warnings
	cfClient, err := broker.getClient()
	if err != nil {
		return ccv3.Package{}, nil, errors.New("Couldn't start session: " + err.Error())
	}

	var pkgCache ccv3.Package

	for pkg.State != constant.PackageReady && pkg.State != constant.PackageFailed && pkg.State != constant.PackageExpired {
		time.Sleep(1000000000)
		ccPkg, warnings, err := cfClient.GetPackage(pkg.GUID)
		broker.Logger.Info("polling package state", lager.Data{
			"package_guid": pkg.GUID,
			"state":        pkg.State,
		})

		broker.showWarnings(warnings, ccPkg)

		allWarnings = append(allWarnings, warnings...)
		if err != nil {
			return ccv3.Package{}, allWarnings, err
		}
		pkgCache = pkg
		pkg = ccv3.Package(ccPkg)
	}

	broker.Logger.Info("polling package final state:", lager.Data{
		"package_guid": pkg.GUID,
		"state":        pkg.State,
	})

	if pkg.State == constant.PackageFailed {
		err := errors.New("package failed")
		broker.Logger.Error(fmt.Sprintf("Service Package Error: Package State %s", pkg.State), err, lager.Data{"Orignal Package": pkgCache, "Checked Package": pkg})
		return ccv3.Package{}, allWarnings, err
	} else if pkg.State == constant.PackageExpired {
		err := errors.New("package expired")
		broker.Logger.Error(fmt.Sprintf("Service Package Error: Package State %s", pkg.State), err, lager.Data{"Orignal Package": pkgCache, "Checked Package": pkg})
		return ccv3.Package{}, allWarnings, err
	}

	return pkg, allWarnings, nil
}

func (broker *ConfigServerBroker) showWarnings(warnings ccv3.Warnings, subject interface{}) {
	broker.Logger.Info(fmt.Sprintf("NOTICE: %d warning(s) were detected!", len(warnings)), lager.Data{"Subject": subject})

	for warn := range warnings {
		w := warnings[warn]
		broker.Logger.Info(fmt.Sprintf("Warning(#%d): %s ", warn+1, w))
	}
}
