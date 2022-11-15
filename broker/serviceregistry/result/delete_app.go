package result

import (
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	brokerapi "github.com/pivotal-cf/brokerapi/v7/domain"
	"github.com/starkandwayne/scs-broker/broker/serviceregistry/registry"
)

type DeleteApp struct {
	Node  *registry.Node
	Spec  brokerapi.DeprovisionServiceSpec
	Error error
}

func NewDeleteApp() *DeleteApp {
	return &DeleteApp{
		Node:  &registry.Node{},
		Spec:  brokerapi.DeprovisionServiceSpec{},
		Error: nil,
	}
}

func (result *DeleteApp) WithApp(app ccv3.Application) *DeleteApp {
	result.Node.App = app

	return result
}

func (result *DeleteApp) WithRoute(route ccv3.Route) *DeleteApp {
	result.Node.Route = route

	return result
}

func (result *DeleteApp) WithNode(node *registry.Node) *DeleteApp {
	result.Node = node

	return result
}

func (result *DeleteApp) WithSpec(spec brokerapi.DeprovisionServiceSpec) *DeleteApp {
	result.Spec = spec

	return result
}

func (result *DeleteApp) WithError(err error) *DeleteApp {
	result.Error = err

	return result
}

func (result *DeleteApp) Failure() error {
	return result.Error
}
