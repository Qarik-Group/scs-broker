package serviceregistry

import (
	"errors"
	"strings"

	"code.cloudfoundry.org/lager"
	"github.com/google/uuid"

	//"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"github.com/starkandwayne/scs-broker/broker/serviceregistry/result"
	"github.com/starkandwayne/scs-broker/broker/serviceregistry/worker"
	"github.com/starkandwayne/scs-broker/broker/utilities"
	"github.com/starkandwayne/scs-broker/logger"
)

type node struct {
	Name      string
	ServiceID string
}

// func (broker *Broker) deployRegistry(space ccv3.Space, serviceId string, desired int) ([]*result.PushApp, error) {
func (broker *Broker) deployRegistry(serviceId string, instanceId string, desired int) (result.PushAppCollection, error) {
	pool := worker.New()
	defer pool.StopWait()

	pushedApps := make([]*result.PushApp, 0)

	pipeline := make(chan *result.PushApp, desired)
	nodes := make(chan *node, desired)

	for _, name := range utilities.NodeNames(instanceId, desired) {
		nodes <- &node{Name: name, ServiceID: serviceId}
	}

	for i := 0; i < desired; i++ {
		jobID, _ := uuid.NewRandom()

		pool.Submit(func() {
			node := <-nodes
			logger.Info(
				"deploying node",
				lager.Data{"job-id": jobID, "node-name": node.Name},
			)

			//broker.deployNode(space, serviceId, nodeName, pipeline)
			broker.deployNode(node.ServiceID, node.Name, pipeline)
		})

	}

	for len(pushedApps) < desired {
		pushedApps = append(pushedApps, <-pipeline)
	}

	errorsPresent := false
	for _, p := range pushedApps {
		if p.Error != nil {
			errorsPresent = true
		}
	}

	if errorsPresent {
		errs := make([]string, 0)
		for _, p := range pushedApps {
			if p.Error != nil {
				errs = append(errs, p.Error.Error())
			}
		}

		return pushedApps, errors.New("errors happened while deploying registry nodes: " + strings.Join(errs, ", "))
	}

	return pushedApps, nil
}
