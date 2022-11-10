package broker

import "code.cloudfoundry.org/lager"

func (broker *SCSBroker) logWorkflowError(msg string, workflow string, err error) {
	broker.Logger.Info(msg, lager.Data{"workflow": workflow, "error": err.Error()})
}

func (broker *SCSBroker) logDeployNodeInfo(msg string, app string, space string) {
	broker.Logger.Info(msg, lager.Data{"app": app, "space": space})
}
