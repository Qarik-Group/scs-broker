package serviceregistry

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/starkandwayne/scs-broker/broker/serviceregistry/registry"
	"github.com/starkandwayne/scs-broker/broker/serviceregistry/result"
	"github.com/starkandwayne/scs-broker/broker/serviceregistry/worker"
	"github.com/starkandwayne/scs-broker/broker/utilities"
	"github.com/starkandwayne/scs-broker/logger"
)

func (broker *Broker) scaleDown(serviceID string, instanceID string, existing []*registry.Node, desiredCount int) ([]*registry.Node, error) {
	if desiredCount < 0 {
		desiredCount = 0
	}

	toDelete := utilities.CalculateDeleteable(existing, desiredCount)

	pool := worker.New()
	defer pool.StopWait()

	deletedApps := make([]*result.DeleteApp, 0)

	pipeline := make(chan *result.DeleteApp, len(toDelete))
	deletionQueue := make(chan *registry.Node, len(toDelete))

	for _, app := range toDelete {
		deletionQueue <- app
	}

	for i := 0; i < len(toDelete); i++ {
		jobID, _ := uuid.NewRandom()

		pool.Submit(func() {
			node := <-deletionQueue
			logger.JobInfo("deleteing node", jobID.String(), node.App.Name)

			broker.deleteNode(node, pipeline)
		})

	}

	for len(deletedApps) < len(toDelete) {
		deletedApps = append(deletedApps, <-pipeline)
	}

	errs := make([]string, 0)
	for _, p := range deletedApps {
		if p.Error != nil {
			errs = append(errs, p.Error.Error())
		}
	}

	if len(errs) > 0 {
		return existing, errors.New("errors happened while deleting registry nodes: " + strings.Join(errs, ", "))
	}

	remainder := existing[:len(existing)-(1+len(toDelete))]

	return remainder, nil
}
