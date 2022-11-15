package logger

import "code.cloudfoundry.org/lager"

func WorkflowError(msg string, workflow string, err error) {
	Error(msg, err, lager.Data{"workflow": workflow})
}
