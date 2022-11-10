package broker

import (
	"errors"
	"strings"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/lager"

	"github.com/google/uuid"
	"github.com/starkandwayne/scs-broker/broker/result"
	"github.com/starkandwayne/scs-broker/broker/utilities"
	"github.com/starkandwayne/scs-broker/broker/worker"
)

type unupdatedNode struct {
	App ccv3.Application
	URL string
}

func (broker *SCSBroker) updateRegistry(deployed []*result.PushApp, rc *utilities.RegistryConfig) ([]*result.UpdateApp, error) {
	pool := worker.New()
	defer pool.StopWait()
	updatedApps := make([]*result.UpdateApp, 0)

	pipeline := make(chan *result.UpdateApp, len(deployed))
	nodes := make(chan *unupdatedNode, len(deployed))

	for _, push := range deployed {
		nodes <- &unupdatedNode{App: push.App, URL: push.Route.URL}
	}

	for i := 0; i < len(deployed); i++ {
		jobID, _ := uuid.NewRandom()

		pool.Submit(func() {
			node := <-nodes
			broker.Logger.Info(
				"updating node",
				lager.Data{"job-id": jobID, "node-name": node.App.Name},
			)

			//broker.deployNode(space, serviceId, nodeName, pipeline)
			broker.updateNode(node, rc, pipeline)
		})
	}

	for len(updatedApps) < len(deployed) {
		updatedApps = append(updatedApps, <-pipeline)
	}

	errorsPresent := false
	for _, p := range updatedApps {
		if p.Error != nil {
			errorsPresent = true
		}
	}

	if errorsPresent {
		errs := make([]string, 0)
		for _, p := range updatedApps {
			if p.Error != nil {
				errs = append(errs, p.Error.Error())
			}
		}

		return updatedApps, errors.New("errors happened while updating registry nodes: " + strings.Join(errs, ", "))
	}

	return updatedApps, nil
}
