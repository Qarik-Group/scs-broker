package broker

import (
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
)

func (broker *SCSBroker) loadRoutes(nodes []ccv3.Application) ([]ccv3.Route, error) {
	filtered := make([]ccv3.Route, 0)

	client, err := broker.GetClient()
	if err != nil {
		return filtered, err
	}

	for _, app := range nodes {
		candidates, _, err := client.GetApplicationRoutes(app.GUID)
		if err != nil {
			return filtered, err
		}

		if len(candidates) > 0 {
			filtered = append(filtered, candidates[0])
		}
	}

	return filtered, nil
}
