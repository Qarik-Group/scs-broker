package serviceregistry

import (
	"context"

	brokerapi "github.com/pivotal-cf/brokerapi/v7/domain"
	"github.com/starkandwayne/scs-broker/broker/utilities"
	"github.com/starkandwayne/scs-broker/config"
)

func (broker *Broker) Update(cxt context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	spec := brokerapi.UpdateServiceSpec{}

	rp, err := utilities.ExtractRegistryParams(string(details.RawParameters))
	if err != nil {
		return spec, err
	}

	desiredCount, err := rp.Count()
	if err != nil {
		return spec, err
	}

	deployed, err := broker.loadNodes(config.Parsed.InstanceSpaceGUID, instanceID)
	if err != nil {
		return spec, err
	}

	return spec, broker.scaleRegistry(details.ServiceID, instanceID, deployed, desiredCount)
}
