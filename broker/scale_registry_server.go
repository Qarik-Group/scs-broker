package broker

import (
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/cli/types"
	"code.cloudfoundry.org/lager"
)

func (broker *SCSBroker) scaleRegistryServer(cfClient *ccv3.Client, app *ccv3.Application, count int) error {
	p := ccv3.Process{
		Type:       "web",
		Instances:  types.NullInt{Value: count, IsSet: true},
		MemoryInMB: types.NullUint64{Value: 0, IsSet: false},
		//DiskInDB:   types.NullUint64{Value: 0, IsSet: false},
	}

	tentative, _, err := cfClient.CreateApplicationProcessScale(app.GUID, p)
	if err != nil {
		broker.Logger.Info("trying to scale the app raised an error", lager.Data{
			"app_guid":     app.GUID,
			"process_guid": p.GUID,
			"error":        err.Error(),
		})
		return err
	}

	_, _, err = broker.pollScale(tentative, count)

	return err
}
