package broker

import (
	"context"
	"fmt"

	brokerapi "github.com/pivotal-cf/brokerapi/domain"
	"github.com/starkandwayne/scs-broker/broker/utilities"
)

func (broker *SCSBroker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails, asyncAllowed bool) (brokerapi.UnbindSpec, error) {
	unbind := brokerapi.UnbindSpec{}

	broker.Logger.Info("UnBind: GetUAAClient")
	api, err := broker.GetUaaClient()
	if err != nil {
		broker.Logger.Info("UnBind: Error in GetUAAClient")
		return unbind, err
	}

	broker.Logger.Info("UnBind: makeClientIdForBinding")
	clientId := utilities.MakeClientIdForBinding(details.ServiceID, bindingID)

	broker.Logger.Info(fmt.Sprintf("UnBind: DeleteClient bindingID:%s clientid %s", bindingID, clientId))
	_, err = api.DeleteClient(clientId)
	if err != nil {
		broker.Logger.Error("UnBind: Error in DeleteClient - will attempt to remove anyway", err)
		return unbind, nil
	}
	broker.Logger.Info("UnBind: Return")
	return unbind, nil
}
