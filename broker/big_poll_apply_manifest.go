package broker

import (
	"errors"
	"fmt"
	"time"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
)

func (broker *SCSBroker) PollApplyManifest(jobURL ccv3.JobURL) (ccv3.Warnings, error) {
	var allWarnings ccv3.Warnings
	cfClient, err := broker.GetClient()
	if err != nil {
		return nil, errors.New("Couldn't start session: " + err.Error())
	}
	count := 0

	for {
		time.Sleep(time.Second)
		count += 1

		broker.Logger.Info(fmt.Sprintf("Polling iteration %d, job %s", count, jobURL))
		job, warnings, err := cfClient.GetJob(jobURL)
		allWarnings = append(allWarnings, warnings...)
		broker.showWarnings(warnings, "poll-apply-manifest")
		if err != nil {
			return allWarnings, err
		}

		broker.Logger.Info(fmt.Sprintf("HERE'S THE FUCKING JOB STATE: %s", job.State))
		if job.HasFailed() {
			err = job.Errors()[0]
			broker.logWorkflowError("pollApplyManifest", "*none*", err)
			return allWarnings, err
		}

		if job.IsComplete() {
			return allWarnings, nil
		}
	}

	return allWarnings, nil
}
