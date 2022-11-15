package logger

import "code.cloudfoundry.org/lager"

func JobInfo(msg string, jobID string, nodeName string) {
	Info(msg, lager.Data{"job-id": jobID, "node-name": nodeName})
}
