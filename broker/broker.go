package broker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3/constant"
	"code.cloudfoundry.org/cli/types"
	"code.cloudfoundry.org/cli/util/configv3"
	brokerapi "github.com/pivotal-cf/brokerapi/domain"
	brokerapiresponses "github.com/pivotal-cf/brokerapi/domain/apiresponses"
	"github.com/starkandwayne/config-server-broker/config"
)

type ConfigServerBroker struct {
	Config config.Config
}

func (broker *ConfigServerBroker) Services(ctx context.Context) ([]brokerapi.Service, error) {
	planList := []brokerapi.ServicePlan{
		brokerapi.ServicePlan{
			ID:          broker.Config.BasicPlanId,
			Name:        broker.Config.BasicPlanName,
			Description: "This plan provides a config server deployed to cf",
			Metadata: &brokerapi.ServicePlanMetadata{
				DisplayName: "Basic",
			},
		}}

	return []brokerapi.Service{
		brokerapi.Service{
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

type InstanceParams struct {
	GitRepoUrl string `json:"gitRepoUrl"`
}

func (broker *ConfigServerBroker) Provision(ctx context.Context, instanceID string, serviceDetails brokerapi.ProvisionDetails, asyncAllowed bool) (spec brokerapi.ProvisionedServiceSpec, err error) {
	spec = brokerapi.ProvisionedServiceSpec{}

	var params InstanceParams
	err = json.Unmarshal(serviceDetails.RawParameters, &params)
	if err != nil {
		return spec, err
	}
	if params.GitRepoUrl == "" {
		return spec, errors.New("Missing parameter 'gitRepoUrl'")
	}

	if serviceDetails.PlanID != broker.Config.BasicPlanId {
		return spec, errors.New("plan_id not recognized")
	}

	err = broker.createBasicInstance(instanceID, params)
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
	_, _, err = cfClient.DeleteRoute(routes[0].GUID)
	_, _, err = cfClient.DeleteApplication(app.GUID)
	if err != nil {
		return spec, err
	}

	return spec, nil
}

func (broker *ConfigServerBroker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails, asyncAllowed bool) (brokerapi.UnbindSpec, error) {
	return brokerapi.UnbindSpec{}, brokerapiresponses.ErrInstanceDoesNotExist
}

func (broker *ConfigServerBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (brokerapi.Binding, error) {
	return brokerapi.Binding{}, brokerapiresponses.ErrInstanceDoesNotExist
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
	return brokerapi.LastOperation{}, errors.New("not implemented")
}

func makeAppName(instanceId string) string {
	return "config-server-" + instanceId
}
func (broker *ConfigServerBroker) createBasicInstance(instanceId string, params InstanceParams) error {
	cfClient, err := broker.getClient()
	if err != nil {
		return errors.New("Couldn't start session: " + err.Error())
	}
	appName := makeAppName(instanceId)
	spaceGUID := broker.Config.InstanceSpaceGUID

	app, _, err := cfClient.CreateApplication(
		ccv3.Application{
			Name:          appName,
			LifecycleType: constant.AppLifecycleTypeDocker,
			State:         constant.ApplicationStarted,
			Relationships: ccv3.Relationships{
				constant.RelationshipTypeSpace: ccv3.Relationship{GUID: spaceGUID},
			},
		},
	)
	if err != nil {
		return err
	}
	pkg, _, err := cfClient.CreatePackage(
		ccv3.Package{
			Type: constant.PackageTypeDocker,
			Relationships: ccv3.Relationships{
				constant.RelationshipTypeApplication: ccv3.Relationship{GUID: app.GUID},
			},
			DockerImage: "hyness/spring-cloud-config-server:latest",
		})
	if err != nil {
		return err
	}
	build, _, err := cfClient.CreateBuild(ccv3.Build{PackageGUID: pkg.GUID})
	if err != nil {
		return err
	}

	droplet, _, err := broker.pollBuild(build.GUID, appName)
	if err != nil {
		return err
	}
	_, _, err = cfClient.SetApplicationDroplet(app.GUID, droplet.GUID)
	if err != nil {
		return err
	}
	_, _, err = cfClient.UpdateApplicationEnvironmentVariables(app.GUID, ccv3.EnvironmentVariables{
		"SPRING_CLOUD_CONFIG_SERVER_GIT_URI": *types.NewFilteredString(params.GitRepoUrl),
	})
	domains, _, err := cfClient.GetDomains(
		ccv3.Query{Key: ccv3.NameFilter, Values: []string{broker.Config.InstanceDomain}},
	)
	if err != nil {
		return err
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
		return ccv3.Droplet{}, nil, errors.New("Couldn't start session: " + err.Error())
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

			interval.Reset(configv3.DefaultStagingTimeout)

		case <-timeout:
			return ccv3.Droplet{}, allWarnings, errors.New("Staging timed out")
		}
	}
}
