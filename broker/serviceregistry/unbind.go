package serviceregistry

import (
	"context"

	brokerapi "github.com/pivotal-cf/brokerapi/v7/domain"
	"github.com/starkandwayne/scs-broker/logger"
)

func (broker *Broker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails, asyncAllowed bool) (brokerapi.UnbindSpec, error) {
	unbind := brokerapi.UnbindSpec{}

	logger.Info("unbindRegistryServer: nothing to clean up")
	return unbind, nil
}
