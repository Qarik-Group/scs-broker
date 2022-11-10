package broker

import (
	"context"

	brokerapi "github.com/pivotal-cf/brokerapi/v7/domain"
)

func (broker *SCSBroker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails, asyncAllowed bool) (brokerapi.UnbindSpec, error) {
	var unbinder func(context.Context, string, string, brokerapi.UnbindDetails, bool) (brokerapi.UnbindSpec, error)

	switch details.ServiceID {
	case "service-registry":
		unbinder = broker.unbindRegistryServerInstance
	case "config-server":
		unbinder = broker.unbindConfigServerInstance
	}

	return unbinder(ctx, instanceID, bindingID, details, asyncAllowed)
}
