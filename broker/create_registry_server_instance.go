package broker

import (
	"fmt"

	"github.com/starkandwayne/scs-broker/broker/utilities"
)

func (broker *SCSBroker) createRegistryServerInstance(serviceId string, instanceId string, jsonparams string, params map[string]string) (string, error) {
	//service, err := broker.GetServiceByServiceID(serviceId)
	//if err != nil {
	//return "", err
	//}

	//client, err := broker.GetClient()
	//if err != nil {
	//return "", err
	//}

	// get target org
	//orgGUID := service.ServiceOrganizationGUID
	//org, _, err := client.GetOrganization(orgGUID)
	//if err != nil {
	//return "", err
	//}

	// create target space
	//space, err := broker.createRegistrySpace(org, instanceId)
	//if err != nil {
	//return "", nil
	//}

	// concurrently create $COUNT apps
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
			rc.AddPeer("https", pushApp.Route.URL)
		}
	} else {
		rc.Standalone()
	}

	_, err = broker.updateRegistry(deployed, rc)
	if err != nil {
		return "", err
	}

	// restart all apps

	//err := broker.restartRegistry(updated)
	//if err != nil {
	//return "", err
	//}

	return fmt.Sprintf("service-registry-%s.%s", instanceId, broker.Config.InstanceDomain), nil
}
