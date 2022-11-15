package broker

import (
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/cli/types"
)

func (broker *SCSBroker) scaleRegistryServer(cfClient *ccv3.Client, app *ccv3.Application, count int) error {
	p := ccv3.Process{
		Type:       "web",
		Instances:  types.NullInt{Value: count, IsSet: true},
		MemoryInMB: types.NullUint64{Value: 0, IsSet: false},
		//DiskInDB:   types.NullUint64{Value: 0, IsSet: false},
	}

	tentative, _, err := cfClient.CreateApplicationProcessScale(app.GUID, p)

	_, _, err = broker.pollScale(tentative, count)

	return err
}
