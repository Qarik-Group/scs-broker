package broker

import (
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	clients "github.com/cloudfoundry-community/go-cf-clients-helper"
	cf "github.com/cloudfoundry-community/go-cfclient"
	"github.com/cloudfoundry-community/go-uaa"
)

func (broker *SCSBroker) GetClient() (*ccv3.Client, error) {

	config := clients.Config{
		Endpoint:          broker.Config.CfConfig.ApiUrl,
		SkipSslValidation: broker.Config.CfConfig.SkipSslValidation,
		User:              broker.Config.CfConfig.CfUsername,
		Password:          broker.Config.CfConfig.CfPassword,
	}

	session, err := clients.NewSession(config)
	if err != nil {
		return nil, err
	}
	return session.V3(), err
}

func (broker *SCSBroker) GetCommunity() (*cf.Client, error) {
	config := &cf.Config{
		ApiAddress:        broker.Config.CfConfig.ApiUrl,
		SkipSslValidation: broker.Config.CfConfig.SkipSslValidation,
		Username:          broker.Config.CfConfig.CfUsername,
		Password:          broker.Config.CfConfig.CfPassword,
	}

	return cf.NewClient(config)
}

func (broker *SCSBroker) GetUaaClient() (*uaa.API, error) {

	cf, err := broker.GetClient()
	if err != nil {
		return nil, err
	}
	info, _, _, err := cf.GetInfo()
	if err != nil {
		return nil, err
	}

	uaaClient, err := uaa.New(info.UAA(), uaa.WithClientCredentials(broker.Config.CfConfig.UaaClientID, broker.Config.CfConfig.UaaClientSecret, uaa.JSONWebToken), uaa.WithSkipSSLValidation(broker.Config.CfConfig.SkipSslValidation))
	if err != nil {
		return nil, err
	}
	return uaaClient, err
}
