package serviceregistry

import (
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"github.com/starkandwayne/scs-broker/broker/serviceregistry/registry"
	"github.com/starkandwayne/scs-broker/client"
)

func (broker *Broker) loadNodes(spaceGUID string, instanceID string) ([]*registry.Node, error) {
	output := make([]*registry.Node, 0)

	client, err := client.GetClient()
	if err != nil {
		return output, err
	}

	candidates, _, err := client.GetApplications(
		ccv3.Query{Key: ccv3.SpaceGUIDFilter, Values: []string{spaceGUID}},
	)
	if err != nil {
		return output, err
	}

	routeErrors := make([]string, 0)

	prefix := "service-registry-" + instanceID + "-"
	for _, prospect := range candidates {
		if strings.HasPrefix(prospect.Name, prefix) {
			routes, _, err := client.GetApplicationRoutes(prospect.GUID)
			if err != nil {
				routeErrors = append(routeErrors, fmt.Sprintf("error loading routes for %s - %s", prospect.Name, err.Error()))
			}
			output = append(output, &registry.Node{App: prospect, Route: routes[0]})
		}
	}

	if len(routeErrors) > 0 {
		err = fmt.Errorf("got errors loading nodes: %s", strings.Join(routeErrors, ", "))
	}

	return output, err
}
