package broker

import (
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"github.com/cloudfoundry-community/go-cf-clients-helper"
	"github.com/cloudfoundry-community/go-uaa"
)

func (broker *ConfigServerBroker) getClient() (*ccv3.Client, error) {

	config := clients.Config{
		Endpoint:          broker.Config.CfConfig.ApiUrl,
		SkipSslValidation: broker.Config.CfConfig.SkipSslValidation,
		User:              broker.Config.CfConfig.Username,
		Password:          broker.Config.CfConfig.Password,
	}

	session, err := clients.NewSession(config)
	if err != nil {
		return nil, err
	}
	return session.V3(), err
}

func (broker *ConfigServerBroker) getUaaClient() (*uaa.API, error) {

	cf, err := broker.getClient()
	if err != nil {
		return nil, err
	}
	info, _, _, err := cf.GetInfo()
	if err != nil {
		return nil, err
	}

	uaaClient, err := uaa.New(info.UAA(), uaa.WithClientCredentials(broker.Config.UaaConfig.ClientID, broker.Config.UaaConfig.ClientSecret, uaa.JSONWebToken), uaa.WithSkipSSLValidation(broker.Config.CfConfig.SkipSslValidation))
	if err != nil {
		return nil, err
	}
	return uaaClient, err
}
