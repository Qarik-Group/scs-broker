package broker

import (
	"context"

	brokerapi "github.com/pivotal-cf/brokerapi/domain"
)

func (broker *ConfigServerBroker) Update(cxt context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	spec := brokerapi.UpdateServiceSpec{}

	kind, err := getKind(details)
	if err != nil {
		return spec, err
	}

	var updater func(context.Context, string, brokerapi.UpdateDetails, bool) (brokerapi.UpdateServiceSpec, error)

	switch kind {
	case "registry-server":
		updater = broker.updateRegistryServerInstance
	case "config-server":
		updater = broker.updateConfigServerInstance
	}

	return updater(cxt, instanceID, details, asyncAllowed)
}
