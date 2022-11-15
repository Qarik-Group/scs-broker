package serviceregistry

import (
	"fmt"

	"github.com/starkandwayne/scs-broker/broker/utilities"
	"github.com/starkandwayne/scs-broker/config"
)

func (broker *Broker) createInstance(serviceId string, instanceId string, jsonparams string, params map[string]string) (string, error) {
	rp, err := utilities.ExtractRegistryParams(string(jsonparams))
	if err != nil {
		return "", err
	}

	desiredCount, err := rp.Count()
	if err != nil {
		return "", err
	}

	rc := utilities.NewRegistryConfig()

	//_, err = broker.deployRegistry(space, serviceId, desiredCount)
	deployed, err := broker.deployRegistry(serviceId, instanceId, desiredCount)
	if err != nil {
		return "", err
	}

	// update all apps with a proper config
	if desiredCount > 1 {
		rc.Clustered()

		for _, pushApp := range deployed {
			rc.AddPeer("https", pushApp.Node.Route.URL)
		}
	} else {
		rc.Standalone()
	}

	_, err = broker.updateRegistry(deployed.Nodes())
	if err != nil {
		return "", err
	}

	// restart all apps

	//err := broker.restartRegistry(updated)
	//if err != nil {
	//return "", err
	//}

	return fmt.Sprintf("service-registry-%s.%s", instanceId, config.Parsed.InstanceDomain), nil
}
