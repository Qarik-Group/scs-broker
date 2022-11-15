package serviceregistry

import (
	"github.com/starkandwayne/scs-broker/broker/serviceregistry/registry"
)

func (broker *Broker) scaleRegistry(serviceID string, instanceID string, existing []*registry.Node, desiredCount int) error {
	currentCount := len(existing)

	var nodes []*registry.Node
	var err error

	//if desiredCount > currentCount {
	//// scale up
	//nodes, err := broker.scaleUp(serviceID, instanceID, existing, desiredCount)
	//if err != nil {
	//return err
	//}
	//}

	if desiredCount < currentCount {
		// scale down
		nodes, err = broker.scaleDown(serviceID, instanceID, existing, desiredCount)
		if err != nil {
			return err
		}
	}

	if len(nodes) > 0 {
		// the node count changed, so we need to update and restart

		_, err = broker.updateRegistry(nodes)
		if err != nil {
			return err
		}

		//_, err = broker.restartRegistry()
		//if err != nil {
		//return err
		//}
	}

	// the desired count of nodes already exists, don't scale
	return nil
}
