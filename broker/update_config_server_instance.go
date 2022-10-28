package broker

import (
	"context"
	"errors"

	"code.cloudfoundry.org/lager"
	brokerapi "github.com/pivotal-cf/brokerapi/domain"
	scsccparser "github.com/starkandwayne/spring-cloud-services-cli-config-parser"
)

func (broker *ConfigServerBroker) updateConfigServerInstance(cxt context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	spec := brokerapi.UpdateServiceSpec{}

	kind, err := getKind(details)
	if err != nil {
		return spec, err
	}

	appName := makeAppName(kind, instanceID)
	spaceGUID := broker.Config.InstanceSpaceGUID

	broker.Logger.Info("update-service-instance", lager.Data{"plan-id": details.PlanID, "service-id": details.ServiceID})
	envsetup := scsccparser.EnvironmentSetup{}
	cfClient, err := broker.getClient()

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

	if details.PlanID != broker.Config.BasicPlanId {
		return spec, errors.New("plan_id not recognized")
	}

	broker.Logger.Info("Updating Environment")
	err = broker.updateAppEnvironment(cfClient, &app, &info, kind, instanceID, string(details.RawParameters), mapparams)
	if err != nil {
		return spec, err
	}

	_, _, err = cfClient.UpdateApplication(app)
	if err != nil {
		return spec, err
	}

	_, _, err = cfClient.UpdateApplicationRestart(app.GUID)
	if err != nil {
		return spec, err
	}

	return spec, nil
}
