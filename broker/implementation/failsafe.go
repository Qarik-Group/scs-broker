package implementation

import (
	"context"
	"fmt"

	"github.com/pivotal-cf/brokerapi/v7/domain"
)

func unimplemented(serviceID string) error {
	return fmt.Errorf("no broker implemented for service %s", serviceID)
}

type failsafe struct{}

func (broker *failsafe) Provision(ctx context.Context, instanceID string, details domain.ProvisionDetails, asyncAllowed bool) (domain.ProvisionedServiceSpec, error) {
	return domain.ProvisionedServiceSpec{}, unimplemented(details.ServiceID)
}

func (broker *failsafe) Deprovision(ctx context.Context, instanceID string, details domain.DeprovisionDetails, asyncAllowed bool) (domain.DeprovisionServiceSpec, error) {
	return domain.DeprovisionServiceSpec{}, unimplemented(details.ServiceID)
}

func (broker *failsafe) Update(ctx context.Context, instanceID string, details domain.UpdateDetails, asyncAllowed bool) (domain.UpdateServiceSpec, error) {
	return domain.UpdateServiceSpec{}, unimplemented(details.ServiceID)
}

func (broker *failsafe) Bind(ctx context.Context, instanceID string, bindingID string, details domain.BindDetails, asyncAllowed bool) (domain.Binding, error) {
	return domain.Binding{}, unimplemented(details.ServiceID)
}

func (broker *failsafe) Unbind(ctx context.Context, instanceID string, bindingID string, details domain.UnbindDetails, asyncAllowed bool) (domain.UnbindSpec, error) {
	return domain.UnbindSpec{}, unimplemented(details.ServiceID)
}
