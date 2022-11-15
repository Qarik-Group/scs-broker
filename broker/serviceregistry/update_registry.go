package serviceregistry

import (
	"errors"
	"strings"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"

	"github.com/google/uuid"
	"github.com/starkandwayne/scs-broker/broker/serviceregistry/registry"
	"github.com/starkandwayne/scs-broker/broker/serviceregistry/result"
	"github.com/starkandwayne/scs-broker/broker/serviceregistry/worker"
	"github.com/starkandwayne/scs-broker/broker/utilities"
	"github.com/starkandwayne/scs-broker/logger"
)

type unupdatedNode struct {
	App ccv3.Application
	URL string
}

func (broker *Broker) updateRegistry(deployed []*registry.Node) ([]*result.UpdateApp, error) {
	rc := utilities.NewRegistryConfig()

	if len(deployed) > 1 {
		rc.Clustered()

		for _, node := range deployed {
			rc.AddPeer("https", node.Route.URL)
		}
	} else {
		rc.Standalone()
	}

	pool := worker.New()
	defer pool.StopWait()
	updatedApps := make([]*result.UpdateApp, 0)

	pipeline := make(chan *result.UpdateApp, len(deployed))
	nodes := make(chan *registry.Node, len(deployed))

	for _, node := range deployed {
		nodes <- node
	}

	//for i := 0; i < len(deployed); i++ {
	for _, _ = range deployed {
		jobID, _ := uuid.NewRandom()

		pool.Submit(func() {
			node := <-nodes
			logger.JobInfo("updating node", jobID.String(), node.App.Name)

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
