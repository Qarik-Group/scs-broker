package configserver

import (
	"context"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccerror"
	brokerapi "github.com/pivotal-cf/brokerapi/v7/domain"
	"github.com/starkandwayne/scs-broker/broker/utilities"
	"github.com/starkandwayne/scs-broker/client"
	"github.com/starkandwayne/scs-broker/config"
	"github.com/starkandwayne/scs-broker/logger"
)

func (broker *Broker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	spec := brokerapi.DeprovisionServiceSpec{}

	cfClient, err := client.GetClient()
	if err != nil {
		return spec, err
	}
	appName := utilities.MakeAppName(details.ServiceID, instanceID)
	app, _, err := cfClient.GetApplicationByNameAndSpace(appName, config.Parsed.InstanceSpaceGUID)
	appNotFound := ccerror.ApplicationNotFoundError{Name: appName}
	if err == appNotFound {
		logger.Info("app-not-found")
		return spec, nil
	}

	if err != nil {
		return spec, err
	}
	routes, _, err := cfClient.GetApplicationRoutes(app.GUID)
	if err != nil {
		return spec, err
	}
	_, _, err = cfClient.UpdateApplicationStop(app.GUID)
	if err != nil {
		return spec, err
	}

	for route := range routes {
		_, _, err := cfClient.DeleteRoute(routes[route].GUID)
		if err != nil {
			return spec, err
		}
	}

	_, _, err = cfClient.DeleteApplication(app.GUID)
	if err != nil {
		return spec, err
	}

	return spec, nil
}
