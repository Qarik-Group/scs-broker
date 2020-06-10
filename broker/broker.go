package broker

import (
	"context"
	"errors"
	"fmt"

	brokerapi "github.com/pivotal-cf/brokerapi/domain"
	brokerapiresponses "github.com/pivotal-cf/brokerapi/domain/apiresponses"
	"github.com/starkandwayne/config-server-broker/config"
)

type ConfigServerBroker struct {
	Config config.Config
}

func (broker *ConfigServerBroker) Services(ctx context.Context) ([]brokerapi.Service, error) {
	planList := []brokerapi.ServicePlan{}
	// for _, plan := range broker.plans() {
	// 	planList = append(planList, *plan)
	// }

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

func (broker *ConfigServerBroker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails, asyncAllowed bool) (brokerapi.UnbindSpec, error) {
	return brokerapi.UnbindSpec{}, brokerapiresponses.ErrInstanceDoesNotExist
}

//Provision ...
func (broker *ConfigServerBroker) Provision(ctx context.Context, instanceID string, serviceDetails brokerapi.ProvisionDetails, asyncAllowed bool) (spec brokerapi.ProvisionedServiceSpec, err error) {
	spec = brokerapi.ProvisionedServiceSpec{}
	return spec, errors.New("not Implemented")
}

func (broker *ConfigServerBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	spec := brokerapi.DeprovisionServiceSpec{}
	return spec, brokerapiresponses.ErrInstanceDoesNotExist
}

func (broker *ConfigServerBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (brokerapi.Binding, error) {
	return brokerapi.Binding{}, brokerapiresponses.ErrInstanceDoesNotExist
}

// type ServiceBroker interface {

// 	// GetInstance fetches information about a service instance
// 	//   GET /v2/service_instances/{instance_id}
// 	GetInstance(ctx context.Context, instanceID string) (GetInstanceDetailsSpec, error)

// 	// Update modifies an existing service instance
// 	//  PATCH /v2/service_instances/{instance_id}
// 	Update(ctx context.Context, instanceID string, details UpdateDetails, asyncAllowed bool) (UpdateServiceSpec, error)

// 	// LastOperation fetches last operation state for a service instance
// 	//   GET /v2/service_instances/{instance_id}/last_operation
// 	LastOperation(ctx context.Context, instanceID string, details PollDetails) (LastOperation, error)

// 	// Bind creates a new service binding
// 	//   PUT /v2/service_instances/{instance_id}/service_bindings/{binding_id}
// 	Bind(ctx context.Context, instanceID, bindingID string, details BindDetails, asyncAllowed bool) (Binding, error)

// 	// Unbind deletes an existing service binding
// 	//   DELETE /v2/service_instances/{instance_id}/service_bindings/{binding_id}
// 	Unbind(ctx context.Context, instanceID, bindingID string, details UnbindDetails, asyncAllowed bool) (UnbindSpec, error)

// 	// GetBinding fetches an existing service binding
// 	//   GET /v2/service_instances/{instance_id}/service_bindings/{binding_id}
// 	GetBinding(ctx context.Context, instanceID, bindingID string) (GetBindingSpec, error)

// 	// LastBindingOperation fetches last operation state for a service binding
// 	//   GET /v2/service_instances/{instance_id}/service_bindings/{binding_id}/last_operation
// 	LastBindingOperation(ctx context.Context, instanceID, bindingID string, details PollDetails) (LastOperation, error)
// }

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
