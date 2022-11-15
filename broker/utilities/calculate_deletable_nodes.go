package utilities

import (
	"github.com/starkandwayne/scs-broker/broker/serviceregistry/registry"
)

func CalculateDeleteable(existing []*registry.Node, desired int) []*registry.Node {
	toDelete := make([]*registry.Node, 0)

	for len(existing)-len(toDelete) > desired {
		toDelete = append(toDelete, existing[len(existing)-(1+len(toDelete))])
	}

	return toDelete
}
