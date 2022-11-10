package broker

import (
	"context"

	brokerapi "github.com/pivotal-cf/brokerapi/v7/domain"
)

func (broker *SCSBroker) deprovisionRegistryServerInstance(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	return brokerapi.DeprovisionServiceSpec{}, broker.deleteRegistry(details.ServiceID, instanceID)
}
