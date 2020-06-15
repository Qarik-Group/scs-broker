package broker

import (
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"github.com/cloudfoundry-community/go-cf-clients-helper"
)

func (broker *ConfigServerBroker) getClient() (*ccv3.Client, error) {

	config := clients.Config{
		Endpoint: broker.Config.CfConfig.ApiUrl,
		User:     broker.Config.CfConfig.Username,
		Password: broker.Config.CfConfig.Password,
	}
	session, err := clients.NewSession(config)
	if err != nil {
		return nil, err
	}
	return session.V3(), err
}
