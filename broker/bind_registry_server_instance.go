package broker

import (
	"context"
	"fmt"
	"strings"

	brokerapi "github.com/pivotal-cf/brokerapi/v7/domain"
)

func (broker *SCSBroker) bindRegistryServerInstance(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (brokerapi.Binding, error) {
	binding := brokerapi.Binding{}

	apps, err := broker.loadNodes(broker.Config.InstanceSpaceGUID, instanceID)
	if err != nil {
		return binding, err
	}

	routes, err := broker.loadRoutes(apps)
	if err != nil {
		return binding, err
	}

	peers := make([]string, 0)
	for _, rte := range routes {
		peers = append(peers, fmt.Sprintf("https://%s", rte.URL))
	}

	url := strings.Join(peers, ",")

	broker.Logger.Info("Bind: Building binding Credentials")
	binding.Credentials = map[string]string{
		"url": url,
		"uri": url,
	}

	broker.Logger.Info("Bind: Return")

	return binding, nil
}
