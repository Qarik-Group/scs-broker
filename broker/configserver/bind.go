package configserver

import (
	"context"
	"fmt"

	"github.com/cloudfoundry-community/go-uaa"
	brokerapi "github.com/pivotal-cf/brokerapi/v7/domain"
	"github.com/starkandwayne/scs-broker/broker/utilities"
	"github.com/starkandwayne/scs-broker/client"
	"github.com/starkandwayne/scs-broker/config"
	"github.com/starkandwayne/scs-broker/logger"
)

func (broker *Broker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (brokerapi.Binding, error) {
	dummy := brokerapi.Binding{}

	logger.Info("Bind: GetUAAClient")

	api, err := client.GetUaaClient()
	if err != nil {
		logger.Info("Bind: Error in getting client")
		return dummy, err
	}

	clientId := utilities.MakeClientIdForBinding(details.ServiceID, bindingID)
	password := utilities.GenClientPassword()

	uaaClient := uaa.Client{
		ClientID:             clientId,
		AuthorizedGrantTypes: []string{"client_credentials"},
		Authorities:          []string{fmt.Sprintf("%s.%v.read", details.ServiceID, instanceID)},
		DisplayName:          clientId,
		ClientSecret:         password,
	}

	logger.Info("Bind: got client info")
	logger.Info("Bind: Create Client")
	_, err = api.CreateClient(uaaClient)
	if err != nil {
		logger.Info("Bind: Error in CreateClient")
		return dummy, err
	}

	logger.Info("Bind: GetClient")
	cfClient, err := client.GetClient()
	if err != nil {
		logger.Info("Bind: Error in GetClient")
		return dummy, err
	}

	logger.Info("Bind: Get Info")
	info, _, _, err := cfClient.GetInfo()
	if err != nil {
		logger.Info("Bind: Error in Get Info")

		return dummy, err
	}

	logger.Info("Bind: GetApplicationByNameAndSpace")

	app, _, err := cfClient.GetApplicationByNameAndSpace(utilities.MakeAppName(details.ServiceID, instanceID), config.Parsed.InstanceSpaceGUID)
	if err != nil {
		logger.Info("Bind: Error in GetApplicationByNameAndSpace")
		return dummy, err
	}

	logger.Info("Bind: GetApplicationRoutes")
	routes, _, err := cfClient.GetApplicationRoutes(app.GUID)
	if err != nil {
		logger.Info("Bind: Error in GetApplicationRoutes")
		return dummy, err
	}

	logger.Info("Bind: Building binding Credentials")

	creds := map[string]string{
		"uri":              fmt.Sprintf("https://%v", routes[0].URL),
		"access_token_uri": fmt.Sprintf("%v/oauth/token", info.UAA()),
		"client_id":        clientId,
		"client_secret":    password,
	}

	logger.Info("Bind: Return")

	return brokerapi.Binding{Credentials: creds}, nil
}
