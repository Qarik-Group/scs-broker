package broker

import "code.cloudfoundry.org/lager"

func (broker *SCSBroker) logWorkflowError(msg string, workflow string, err error) {
	broker.Logger.Info(msg, lager.Data{"workflow": workflow, "error": err.Error()})
}
