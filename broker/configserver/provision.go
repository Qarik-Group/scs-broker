package configserver

import (
	"context"
	"fmt"

	brokerapi "github.com/pivotal-cf/brokerapi/v7/domain"
	"github.com/starkandwayne/scs-broker/logger"
	scsccparser "github.com/starkandwayne/spring-cloud-services-cli-config-parser"
)

func (broker *Broker) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (spec brokerapi.ProvisionedServiceSpec, err error) {
	logger.Info(fmt.Sprintf("Got these details: %s", details))
	envsetup := scsccparser.EnvironmentSetup{}

	raw := details.RawParameters
	if len(raw) == 0 {
		raw = []byte("{}")
	}

	mapparams, err := envsetup.ParseEnvironmentFromRaw(raw)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, err
	}

	logger.Info("Provisioning a " + details.ServiceID + " service instance")

	url, err := broker.createInstance(details.ServiceID, instanceID, string(details.RawParameters), mapparams)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, err
	}

	return brokerapi.ProvisionedServiceSpec{DashboardURL: url}, nil
}
