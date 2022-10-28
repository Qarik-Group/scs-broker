package broker

import (
	"context"
	"fmt"

	"code.cloudfoundry.org/lager"
	brokerapi "github.com/pivotal-cf/brokerapi/domain"
	"github.com/starkandwayne/scs-broker/config"
)

const (
	ArtifactsDir string = "artifacts"
)

type ConfigServerBroker struct {
	Config config.Config
	Logger lager.Logger
}

func (broker *ConfigServerBroker) Services(ctx context.Context) ([]brokerapi.Service, error) {
	planList := []brokerapi.ServicePlan{
		{
			ID:          "default",
			Name:        "default",
			Description: "This plan provides SCS servers deployed to cf",
			Metadata: &brokerapi.ServicePlanMetadata{
				DisplayName: "Default",
			},
		}}

	configServer := brokerapi.Service{
		//ID:          broker.Config.ServiceID,
		//Name:        broker.Config.ServiceName,
		//Description: broker.Config.Description,
		ID:          "config-server",
		Name:        "config-server",
		Description: "Broker to create config-servers",
		Bindable:    true,
		Plans:       planList,
		Metadata: &brokerapi.ServiceMetadata{
			//DisplayName:         broker.Config.DisplayName,
			DisplayName:         "config-server",
			LongDescription:     broker.Config.LongDescription,
			DocumentationUrl:    broker.Config.DocumentationURL,
			SupportUrl:          broker.Config.SupportURL,
			ImageUrl:            fmt.Sprintf("data:image/png;base64,%s", broker.Config.IconImage),
			ProviderDisplayName: broker.Config.ProviderDisplayName,
		},
		Tags: []string{
			"snw",
			"config-server",
		},
	}

	registryServer := brokerapi.Service{
		ID:          "registry-server",
		Name:        "registry-server",
		Description: "Broker to create registry-servers",
		Bindable:    true,
		Plans: []brokerapi.ServicePlan{
		{
			ID:          "basic",
			Name:        "basic",
			Description: "This plan provides SCS registry servers deployed to cf",
			Metadata: &brokerapi.ServicePlanMetadata{
				DisplayName: "Basic",
			},
		}},
		Metadata: &brokerapi.ServiceMetadata{
			DisplayName: "registry-server",
			ImageUrl:    fmt.Sprintf("data:image/png;base64,%s", broker.Config.IconImage),
		},
		Tags: []string{
			"snw",
			"registry-server",
		},
	}

	return []brokerapi.Service{configServer, registryServer}, nil
}
