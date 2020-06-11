package broker

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudfoundry-community/go-cfclient"
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

func (broker *ConfigServerBroker) Provision(ctx context.Context, instanceID string, serviceDetails brokerapi.ProvisionDetails, asyncAllowed bool) (spec brokerapi.ProvisionedServiceSpec, err error) {
	spec = brokerapi.ProvisionedServiceSpec{}

	if serviceDetails.PlanID != broker.Config.BasicPlanId {
		return spec, errors.New("plan_id not recognized")
	}

	err = broker.createBasicInstance(instanceID)
	if err != nil {
		return spec, err
	}

	return spec, nil
}

func (broker *ConfigServerBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	spec := brokerapi.DeprovisionServiceSpec{}
	return spec, brokerapiresponses.ErrInstanceDoesNotExist
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

func (broker *ConfigServerBroker) createBasicInstance(instanceId string) error {
	request := cfclient.AppCreateRequest{
		Name:        "config-server-" + instanceId,
		DockerImage: "hyness/spring-cloud-config-server:latest",
		SpaceGuid:   "a7cc4fc8-9161-423c-a7a4-19fab3e1b64d",
		State:       cfclient.APP_STARTED,
		Environment: map[string]interface{}{
			"SPRING_CLOUD_CONFIG_SERVER_GIT_URI": "https://github.com/spring-cloud-samples/config-repo",
		},
	}
	fmt.Println("Creating app: %v", request)
	client, err := broker.Config.GetCfClient()
	if err != nil {
		return errors.New("Couldn't create CF client: " + err.Error())
	}
	_, err = client.CreateApp(request)
	if err != nil {
		return errors.New("Couldn't create app: " + err.Error())
	}
	return nil
}
