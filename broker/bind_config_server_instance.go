package broker

import (
	"context"
	"fmt"

	"github.com/cloudfoundry-community/go-uaa"
	brokerapi "github.com/pivotal-cf/brokerapi/v7/domain"
	"github.com/starkandwayne/scs-broker/broker/utilities"
)

func (broker *SCSBroker) bindConfigServerInstance(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (brokerapi.Binding, error) {
	binding := brokerapi.Binding{}

	broker.Logger.Info("Bind: GetUAAClient")

	api, err := broker.GetUaaClient()
	if err != nil {
		broker.Logger.Info("Bind: Error in getting client")
		return binding, err
	}

	clientId := utilities.MakeClientIdForBinding(details.ServiceID, bindingID)
	password := utilities.GenClientPassword()

	client := uaa.Client{
		ClientID:             clientId,
		AuthorizedGrantTypes: []string{"client_credentials"},
		Authorities:          []string{fmt.Sprintf("%s.%v.read", details.ServiceID, instanceID)},
		DisplayName:          clientId,
		ClientSecret:         password,
	}

	broker.Logger.Info("Bind: got client info")
	broker.Logger.Info("Bind: Create Client")
	_, err = api.CreateClient(client)
	if err != nil {
		broker.Logger.Info("Bind: Error in CreateClient")
		return binding, err
	}

	broker.Logger.Info("Bind: GetClient")
	cfClient, err := broker.GetClient()
	if err != nil {
		broker.Logger.Info("Bind: Error in GetClient")
		return binding, err
	}

	broker.Logger.Info("Bind: Get Info")
	info, _, _, err := cfClient.GetInfo()
	if err != nil {
		broker.Logger.Info("Bind: Error in Get Info")

		return binding, err
	}

	broker.Logger.Info("Bind: GetApplicationByNameAndSpace")

	app, _, err := cfClient.GetApplicationByNameAndSpace(utilities.MakeAppName(details.ServiceID, instanceID), broker.Config.InstanceSpaceGUID)
	if err != nil {
		broker.Logger.Info("Bind: Error in GetApplicationByNameAndSpace")
		return binding, err
	}

	broker.Logger.Info("Bind: GetApplicationRoutes")
	routes, _, err := cfClient.GetApplicationRoutes(app.GUID)
	if err != nil {
		broker.Logger.Info("Bind: Error in GetApplicationRoutes")
		return binding, err
	}

	broker.Logger.Info("Bind: Building binding Credentials")
	binding.Credentials = map[string]string{
		"uri":              fmt.Sprintf("https://%v", routes[0].URL),
		"access_token_uri": fmt.Sprintf("%v/oauth/token", info.UAA()),
		"client_id":        clientId,
		"client_secret":    password,
	}

	broker.Logger.Info("Bind: Return")

	return binding, nil
}
