package result

import (
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"github.com/starkandwayne/scs-broker/broker/serviceregistry/registry"
)

type PushApp struct {
	Node  *registry.Node
	Error error
}

func NewPushApp() *PushApp {
	return &PushApp{
		Node:  &registry.Node{},
		Error: nil,
	}
}

func (result *PushApp) WithApp(app ccv3.Application) *PushApp {
	result.Node.App = app

	return result
}

func (result *PushApp) WithRoute(rte ccv3.Route) *PushApp {
	result.Node.Route = rte

	return result
}

func (result *PushApp) WithNode(node *registry.Node) *PushApp {
	result.Node = node

	return result
}

func (result *PushApp) WithError(err error) *PushApp {
	result.Error = err

	return result
}

func (result *PushApp) Failure() error {
	return result.Error
}

type PushAppCollection []*PushApp

func (c PushAppCollection) Nodes() []*registry.Node {
	nodes := make([]*registry.Node, 0)

	for _, p := range c {
		nodes = append(nodes, p.Node)
	}

	return nodes
}
