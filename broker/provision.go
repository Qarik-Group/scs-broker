package broker

import (
	"context"
	"fmt"

	brokerapi "github.com/pivotal-cf/brokerapi/domain"
)

func (broker *SCSBroker) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (spec brokerapi.ProvisionedServiceSpec, err error) {
	broker.Logger.Info(fmt.Sprintf("Got these details: %s", details))
	spec = brokerapi.ProvisionedServiceSpec{}

	broker.Logger.Info("Provisioning a " + details.ServiceID + " service instance")

	broker.CreateServiceInstances(ctx, instanceID, details, true)

	return spec, nil
}
