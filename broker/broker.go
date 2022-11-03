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

type SCSBroker struct {
	Config config.Config
	Logger lager.Logger
}

func (broker *SCSBroker) GetServiceByServiceID(serviceID string) (config.Service, error) {
	for _, service := range broker.Config.Services {
		if service.ServiceID == serviceID {
			return service, nil
		}
	}

	return config.Service{}, fmt.Errorf("No valid service found for %s", serviceID)
}

func (broker *SCSBroker) Services(ctx context.Context) ([]brokerapi.Service, error) {

	services := []brokerapi.Service{}

	for _, service := range broker.Config.Services {
		brokerService := brokerapi.Service{
			ID:          service.ServiceID,
			Name:        service.ServiceName,
			Description: service.ServiceDescription,
			Bindable:    true,
			Plans: []brokerapi.ServicePlan{
				{
					ID:          service.ServicePlanID,
					Name:        service.ServicePlanName,
					Description: service.ServiceDescription,
					Metadata: &brokerapi.ServicePlanMetadata{
						DisplayName: service.ServicePlanName,
					},
				}},
			Metadata: &brokerapi.ServiceMetadata{
				DisplayName: service.ServiceName,
				ImageUrl:    fmt.Sprintf("data:image/png;base64,%s", broker.Config.IconImage),
			},
			Tags: []string{
				"snw",
				service.ServiceName,
			},
		}
		services = append(services, brokerService)
	}

	return services, nil

}
