package broker

import (
	"context"
	"errors"
	"fmt"

	brokerapi "github.com/pivotal-cf/brokerapi/v7/domain"
	"github.com/starkandwayne/scs-broker/broker/implementation"
	"github.com/starkandwayne/scs-broker/config"
)

const (
	ArtifactsDir string = "artifacts"
)

type SCSBroker struct {
}

func New() *SCSBroker {
	return &SCSBroker{}
}

func (broker *SCSBroker) Services(ctx context.Context) ([]brokerapi.Service, error) {

	services := []brokerapi.Service{}

	for _, service := range config.Parsed.Services {
		brokerService := brokerapi.Service{
			ID:          service.ServiceID,
			Name:        service.ServiceName,
			Description: service.ServiceDescription,
			Bindable:    true,
			Plans: []brokerapi.ServicePlan{
				{
					ID:          service.ServicePlanID,
					Name:        service.ServicePlanName,
					Description: service.ServiceDescription,
					Metadata: &brokerapi.ServicePlanMetadata{
						DisplayName: service.ServicePlanName,
					},
				}},
			Metadata: &brokerapi.ServiceMetadata{
				DisplayName: service.ServiceName,
				ImageUrl:    fmt.Sprintf("data:image/png;base64,%s", config.Parsed.IconImage),
			},
			Tags: []string{
				"snw",
				service.ServiceName,
			},
		}
		services = append(services, brokerService)
	}

	return services, nil
}

func (broker *SCSBroker) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (spec brokerapi.ProvisionedServiceSpec, err error) {
	return implementation.
		ByServiceID(details.ServiceID).
		Provision(ctx, instanceID, details, asyncAllowed)
}

func (broker *SCSBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	return implementation.
		ByServiceID(details.ServiceID).
		Deprovision(ctx, instanceID, details, asyncAllowed)
}

func (broker *SCSBroker) Update(ctx context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	return implementation.
		ByServiceID(details.ServiceID).
		Update(ctx, instanceID, details, asyncAllowed)
}

func (broker *SCSBroker) Bind(ctx context.Context, instanceID string, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (brokerapi.Binding, error) {
	return implementation.
		ByServiceID(details.ServiceID).
		Bind(ctx, instanceID, bindingID, details, asyncAllowed)
}

func (broker *SCSBroker) Unbind(ctx context.Context, instanceID string, bindingID string, details brokerapi.UnbindDetails, asyncAllowed bool) (brokerapi.UnbindSpec, error) {
	return implementation.
		ByServiceID(details.ServiceID).
		Unbind(ctx, instanceID, bindingID, details, asyncAllowed)
}

// Here be unimplemented dragons

func (broker *SCSBroker) LastOperation(ctx context.Context, instanceID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, errors.New("not implemented")
}

func (broker *SCSBroker) GetBinding(ctx context.Context, instanceID, bindingID string) (brokerapi.GetBindingSpec, error) {
	return brokerapi.GetBindingSpec{}, errors.New("not implemented")
}

func (broker *SCSBroker) GetInstance(ctx context.Context, instanceID string) (brokerapi.GetInstanceDetailsSpec, error) {
	return brokerapi.GetInstanceDetailsSpec{}, errors.New("not implemented")
}

func (broker *SCSBroker) LastBindingOperation(ctx context.Context, instanceID, bindingID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, errors.New("not implemented")
}
