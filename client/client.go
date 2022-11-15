package client

import (
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	clients "github.com/cloudfoundry-community/go-cf-clients-helper"
	cf "github.com/cloudfoundry-community/go-cfclient"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/starkandwayne/scs-broker/config"
)

func GetClient() (*ccv3.Client, error) {

	config := clients.Config{
		Endpoint:          config.Parsed.CfConfig.ApiUrl,
		SkipSslValidation: config.Parsed.CfConfig.SkipSslValidation,
		User:              config.Parsed.CfConfig.CfUsername,
		Password:          config.Parsed.CfConfig.CfPassword,
	}

	session, err := clients.NewSession(config)
	if err != nil {
		return nil, err
	}
	return session.V3(), err
}

func GetCommunity() (*cf.Client, error) {
	config := &cf.Config{
		ApiAddress:        config.Parsed.CfConfig.ApiUrl,
		SkipSslValidation: config.Parsed.CfConfig.SkipSslValidation,
		Username:          config.Parsed.CfConfig.CfUsername,
		Password:          config.Parsed.CfConfig.CfPassword,
	}

	return cf.NewClient(config)
}

func GetUaaClient() (*uaa.API, error) {

	cf, err := GetClient()
	if err != nil {
		return nil, err
	}
	info, _, _, err := cf.GetInfo()
	if err != nil {
		return nil, err
	}

	uaaClient, err := uaa.New(info.UAA(), uaa.WithClientCredentials(config.Parsed.CfConfig.UaaClientID, config.Parsed.CfConfig.UaaClientSecret, uaa.JSONWebToken), uaa.WithSkipSSLValidation(config.Parsed.CfConfig.SkipSslValidation))
	if err != nil {
		return nil, err
	}
	return uaaClient, err
}
