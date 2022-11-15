package serviceregistry

import (
	"context"

	brokerapi "github.com/pivotal-cf/brokerapi/v7/domain"
	"github.com/starkandwayne/scs-broker/config"
)

func (broker *Broker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	spec := brokerapi.DeprovisionServiceSpec{}

	deployed, err := broker.loadNodes(config.Parsed.InstanceSpaceGUID, instanceID)
	if err != nil {
		return spec, err
	}

	// so, we could do a whole thing where we explicitly delete all of the
	// nodes ... or we could just scale to zero
	return spec, broker.scaleRegistry(details.ServiceID, instanceID, deployed, 0)
}
