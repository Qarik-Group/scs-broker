package broker

import (
	"context"

	brokerapi "github.com/pivotal-cf/brokerapi/v7/domain"
)

func (broker *SCSBroker) Update(cxt context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {

	var updater func(context.Context, string, brokerapi.UpdateDetails, bool) (brokerapi.UpdateServiceSpec, error)

	switch details.ServiceID {
	case "service-registry":
		updater = broker.updateRegistryServerInstance
	case "config-server":
		updater = broker.updateConfigServerInstance
	}

	return updater(cxt, instanceID, details, asyncAllowed)
}
