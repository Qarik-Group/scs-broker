package broker

import (
	"errors"
	"strings"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/lager"
	"github.com/google/uuid"
	"github.com/starkandwayne/scs-broker/broker/result"
	"github.com/starkandwayne/scs-broker/broker/worker"
)

type deletable struct {
	App ccv3.Application
}

func (broker *SCSBroker) deleteRegistry(serviceID string, instanceID string) error {
	pool := worker.New()
	defer pool.StopWait()

	// load nodes
	apps, err := broker.loadNodes(broker.Config.InstanceSpaceGUID, instanceID)
	if err != nil {
		return err
	}

	deletedApps := make([]*result.DeleteApp, 0)

	pipeline := make(chan *result.DeleteApp, len(apps))
	deletionQueue := make(chan *deletable, len(apps))

	for _, app := range apps {
		deletionQueue <- &deletable{App: app}
	}

	for i := 0; i < len(apps); i++ {
		jobID, _ := uuid.NewRandom()

		pool.Submit(func() {
			node := <-deletionQueue
			broker.Logger.Info(
				"deleting node",
				lager.Data{"job-id": jobID, "node-name": node.App.Name},
			)

			//broker.deployNode(space, serviceID, nodeName, pipeline)
			broker.deleteNode(node, pipeline)
		})

	}

	for len(deletedApps) < len(apps) {
		deletedApps = append(deletedApps, <-pipeline)
	}

	errs := make([]string, 0)
	for _, p := range deletedApps {
		if p.Error != nil {
			errs = append(errs, p.Error.Error())
		}
	}

	if len(errs) > 0 {
		return errors.New("errors happened while deleting registry nodes: " + strings.Join(errs, ", "))
	}

	return nil
}
