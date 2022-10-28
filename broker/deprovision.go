package broker

import (
	"context"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccerror"
	brokerapi "github.com/pivotal-cf/brokerapi/domain"
)

func (broker *ConfigServerBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	spec := brokerapi.DeprovisionServiceSpec{}
	kind, err := getKind(details)
	if err != nil {
		return spec, err
	}

	cfClient, err := broker.getClient()
	if err != nil {
		return spec, err
	}
	appName := makeAppName(kind, instanceID)
	app, _, err := cfClient.GetApplicationByNameAndSpace(appName, broker.Config.InstanceSpaceGUID)
	appNotFound := ccerror.ApplicationNotFoundError{Name: appName}
	if err == appNotFound {
		broker.Logger.Info("app-not-found")
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
