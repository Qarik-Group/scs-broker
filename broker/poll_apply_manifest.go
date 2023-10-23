package broker

import (
	"time"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
)

//func (broker *SCSBroker) pollApplyManifest(jobURL ccv3.JobURL) (ccv3.Warnings, error) {
//var allWarnings ccv3.Warnings
//cfClient, err := broker.GetClient()
//if err != nil {
//return nil, errors.New("Couldn't start session: " + err.Error())
//}
//count := 0

//for {
//time.Sleep(time.Second)
//count += 1

//broker.Logger.Info(fmt.Sprintf("Polling iteration %d, job %s", count, jobURL))
//job, warnings, err := cfClient.GetJob(jobURL)
//allWarnings = append(allWarnings, warnings...)
//broker.showWarnings(warnings, "poll-apply-manifest")
//if err != nil {
//return allWarnings, err
//}

//broker.Logger.Info(fmt.Sprintf("HERE'S THE FUCKING JOB STATE: %s", job.State))
//if job.HasFailed() {
//err = job.Errors()[0]
//broker.logWorkflowError("pollApplyManifest", "*none*", err)
//return allWarnings, err
//}

//if job.IsComplete() {
//return allWarnings, nil
//}
//}

//return allWarnings, nil
//}

func (broker *SCSBroker) pollApplyManifest(jobURL ccv3.JobURL) (ccv3.Warnings, error) {
	var (
		err         error
		warnings    ccv3.Warnings
		allWarnings ccv3.Warnings
		job         ccv3.Job
	)

	client, err := broker.GetClient()
	if err != nil {
		return allWarnings, err
	}

	for {
		job, warnings, err = client.GetJob(jobURL)
		allWarnings = append(allWarnings, warnings...)
		if err != nil {
			return allWarnings, err
		}

		if job.HasFailed() {
			firstError := job.Errors()[0]
			return allWarnings, firstError
		}

		if job.IsComplete() {
			return allWarnings, nil
		}

		time.Sleep(5 * time.Second)
	}
}
