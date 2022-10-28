package broker

import (
	"context"
	"errors"
	"fmt"

	brokerapi "github.com/pivotal-cf/brokerapi/domain"
	scsccparser "github.com/starkandwayne/spring-cloud-services-cli-config-parser"
)

func (broker *ConfigServerBroker) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (spec brokerapi.ProvisionedServiceSpec, err error) {
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

	if details.PlanID != "basic" {
		return spec, errors.New("plan_id not recognized")
	}

	kind, err := getKind(details)
	if err != nil {
		return spec, err
	}

	broker.Logger.Info("Provisioning a " + kind + " service instance")

	var provisioner func(string, string, string, map[string]string) (string, error)

	switch kind {
	case "registry-server":
		provisioner = broker.createRegistryServerInstance
	case "config-server":
		provisioner = broker.createConfigServerInstance

	}

	url, err := provisioner(kind, instanceID, string(details.RawParameters), mapparams)
	if err != nil {
		return spec, err
	}

	spec.DashboardURL = url
	return spec, nil
}
