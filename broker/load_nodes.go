package broker

import (
	"strings"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
)

func (broker *SCSBroker) loadNodes(spaceGUID string, instanceID string) ([]ccv3.Application, error) {
	filtered := make([]ccv3.Application, 0)

	client, err := broker.GetClient()
	if err != nil {
		return filtered, err
	}

	candidates, _, err := client.GetApplications(
		ccv3.Query{Key: ccv3.SpaceGUIDFilter, Values: []string{spaceGUID}},
	)
	if err != nil {
		return filtered, err
	}

	prefix := "service-registry-" + instanceID + "-"
	for _, prospect := range candidates {
		if strings.HasPrefix(prospect.Name, prefix) {
			filtered = append(filtered, prospect)
		}
	}

	return filtered, nil
}
