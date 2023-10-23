package broker

import (
	"context"
	"fmt"

	brokerapi "github.com/pivotal-cf/brokerapi/v7/domain"
	scsccparser "github.com/starkandwayne/spring-cloud-services-cli-config-parser"
)

func (broker *SCSBroker) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (spec brokerapi.ProvisionedServiceSpec, err error) {
	broker.Logger.Info(fmt.Sprintf("Got these details: %s", details))
	spec = brokerapi.ProvisionedServiceSpec{}
	envsetup := scsccparser.EnvironmentSetup{}
	raw := details.RawParameters
	if len(raw) == 0 {
		raw = []byte("{}")
	}
	mapparams, err := envsetup.ParseEnvironmentFromRaw(raw)
	if err != nil {
		return spec, err
	}

	broker.Logger.Info("Provisioning a " + details.ServiceID + " service instance")

	var provisioner func(string, string, string, map[string]string) (string, error)

	switch details.ServiceID {
	case "service-registry":
		provisioner = broker.createRegistryServerInstance
	case "config-server":
		provisioner = broker.createConfigServerInstance

	}

	url, err := provisioner(details.ServiceID, instanceID, string(details.RawParameters), mapparams)
	if err != nil {
		return spec, err
	}

	spec.DashboardURL = url
	return spec, nil
}
