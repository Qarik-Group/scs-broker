package broker

import (
	"errors"
	"time"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3/constant"
	"code.cloudfoundry.org/lager"
)

func (broker *SCSBroker) pollScale(proc ccv3.Process, desired int) (ccv3.Process, ccv3.Warnings, error) {
	var allWarnings ccv3.Warnings
	cfClient, err := broker.GetClient()
	if err != nil {
		return ccv3.Process{}, nil, errors.New("Couldn't start session: " + err.Error())
	}

	done := false

	for !done {
		time.Sleep(1000000000)
		ready := 0
		broker.Logger.Info("polling process instance states", lager.Data{
			"process_guid": proc.GUID,
		})

		instances, warnings, err := cfClient.GetProcessInstances(proc.GUID)
		broker.showWarnings(warnings, proc)
		allWarnings = append(allWarnings, warnings...)
		if err != nil {
			return ccv3.Process{}, allWarnings, err
		}

		for _, instance := range instances {
			if instance.State == constant.ProcessInstanceRunning || instance.State == constant.ProcessInstanceCrashed {
				ready += 1
			}
		}

		if ready == desired {
			done = true
		}
	}

	return proc, allWarnings, nil
}
