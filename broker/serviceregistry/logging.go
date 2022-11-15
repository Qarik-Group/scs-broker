package serviceregistry

import (
	"code.cloudfoundry.org/lager"
	"github.com/starkandwayne/scs-broker/logger"
)

func (broker *Broker) logWorkflowError(msg string, workflow string, err error) {
	logger.Info(msg, lager.Data{"workflow": workflow, "error": err.Error()})
}

func (broker *Broker) logDeployNodeInfo(msg string, app string, space string) {
	logger.Info(msg, lager.Data{"app": app, "space": space})
}
