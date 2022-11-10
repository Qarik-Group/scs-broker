package broker

import (
	"context"

	brokerapi "github.com/pivotal-cf/brokerapi/v7/domain"
)

func (broker *SCSBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	var deprovisioner func(context.Context, string, brokerapi.DeprovisionDetails, bool) (brokerapi.DeprovisionServiceSpec, error)

	switch details.ServiceID {
	case "service-registry":
		deprovisioner = broker.deprovisionRegistryServerInstance
	case "config-server":
		deprovisioner = broker.deprovisionConfigServerInstance
	}

	return deprovisioner(ctx, instanceID, details, asyncAllowed)
}
