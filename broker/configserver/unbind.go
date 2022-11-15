package configserver

import (
	"context"
	"fmt"

	brokerapi "github.com/pivotal-cf/brokerapi/v7/domain"
	"github.com/starkandwayne/scs-broker/broker/utilities"
	"github.com/starkandwayne/scs-broker/client"
	"github.com/starkandwayne/scs-broker/logger"
)

func (broker *Broker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails, asyncAllowed bool) (brokerapi.UnbindSpec, error) {
	unbind := brokerapi.UnbindSpec{}

	logger.Info("UnBind: GetUAAClient")
	api, err := client.GetUaaClient()
	if err != nil {
		logger.Info("UnBind: Error in GetUAAClient")
		return unbind, err
	}

	logger.Info("UnBind: makeClientIdForBinding")
	clientId := utilities.MakeClientIdForBinding(details.ServiceID, bindingID)

	logger.Info(fmt.Sprintf("UnBind: DeleteClient bindingID:%s clientid %s", bindingID, clientId))
	_, err = api.DeleteClient(clientId)
	if err != nil {
		logger.Error("UnBind: Error in DeleteClient - will attempt to remove anyway", err)
		return unbind, nil
	}
	logger.Info("UnBind: Return")
	return unbind, nil
}
