package broker

import (
	"context"

	brokerapi "github.com/pivotal-cf/brokerapi/v7/domain"
)

func (broker *SCSBroker) unbindRegistryServerInstance(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails, asyncAllowed bool) (brokerapi.UnbindSpec, error) {
	unbind := brokerapi.UnbindSpec{}

	broker.Logger.Info("unbindRegistryServer: nothing to clean up")
	return unbind, nil
}
