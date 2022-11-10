package broker

import (
	"context"

	brokerapi "github.com/pivotal-cf/brokerapi/v7/domain"
)

func (broker *SCSBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (brokerapi.Binding, error) {
	var binder func(context.Context, string, string, brokerapi.BindDetails, bool) (brokerapi.Binding, error)

	switch details.ServiceID {
	case "service-registry":
		binder = broker.bindRegistryServerInstance
	case "config-server":
		binder = broker.bindConfigServerInstance
	}

	return binder(ctx, instanceID, bindingID, details, asyncAllowed)
}
