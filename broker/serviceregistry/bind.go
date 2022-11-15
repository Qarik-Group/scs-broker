package serviceregistry

import (
	"context"
	"fmt"
	"strings"

	brokerapi "github.com/pivotal-cf/brokerapi/v7/domain"
	"github.com/starkandwayne/scs-broker/config"
	"github.com/starkandwayne/scs-broker/logger"
)

func (broker *Broker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (brokerapi.Binding, error) {
	dummy := brokerapi.Binding{}

	apps, err := broker.loadNodes(config.Parsed.InstanceSpaceGUID, instanceID)
	if err != nil {
		return dummy, err
	}

	peers := make([]string, 0)
	for _, node := range apps {
		peers = append(peers, fmt.Sprintf("https://%s", node.Route.URL))
	}

	url := strings.Join(peers, ",")

	logger.Info("Bind: Building binding Credentials")
	creds := map[string]string{
		"url": url,
		"uri": url,
	}

	logger.Info("Bind: Return")

	return brokerapi.Binding{Credentials: creds}, nil
}
