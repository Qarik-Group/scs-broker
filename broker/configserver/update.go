package configserver

import (
	"context"
	"errors"

	"code.cloudfoundry.org/lager"
	brokerapi "github.com/pivotal-cf/brokerapi/v7/domain"
	"github.com/starkandwayne/scs-broker/broker/utilities"
	"github.com/starkandwayne/scs-broker/client"
	"github.com/starkandwayne/scs-broker/config"
	"github.com/starkandwayne/scs-broker/logger"
	scsccparser "github.com/starkandwayne/spring-cloud-services-cli-config-parser"
)

func (broker *Broker) Update(cxt context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	spec := brokerapi.UpdateServiceSpec{}

	appName := utilities.MakeAppName(details.ServiceID, instanceID)
	spaceGUID := config.Parsed.InstanceSpaceGUID

	logger.Info("update-service-instance", lager.Data{"plan-id": details.PlanID, "service-id": details.ServiceID})
	envsetup := scsccparser.EnvironmentSetup{}
	cfClient, err := client.GetClient()

	if err != nil {
		return spec, errors.New("Couldn't start session: " + err.Error())
	}

	info, _, _, err := cfClient.GetInfo()
	if err != nil {
		return spec, err
	}

	app, _, err := cfClient.GetApplicationByNameAndSpace(appName, spaceGUID)
	if err != nil {
		return spec, errors.New("Couldn't find app session: " + err.Error())
	}

	mapparams, err := envsetup.ParseEnvironmentFromRaw(details.RawParameters)
	if err != nil {
		return spec, err
	}

	logger.Info("Updating Environment")
	err = broker.UpdateAppEnvironment(cfClient, &app, &info, details.ServiceID, instanceID, string(details.RawParameters), mapparams)
	if err != nil {
		return spec, err
	}

	_, _, err = cfClient.UpdateApplication(utilities.SafeApp(app))
	if err != nil {
		return spec, err
	}

	_, _, err = cfClient.UpdateApplicationRestart(app.GUID)
	if err != nil {
		return spec, err
	}

	return spec, nil
}
