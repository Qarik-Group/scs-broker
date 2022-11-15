package serviceregistry

import (
	brokerapi "github.com/pivotal-cf/brokerapi/v7/domain"
	"github.com/starkandwayne/scs-broker/broker/serviceregistry/registry"
	"github.com/starkandwayne/scs-broker/broker/serviceregistry/result"
	"github.com/starkandwayne/scs-broker/client"
)

func (broker *Broker) deleteNode(node *registry.Node, pipeline chan<- *result.DeleteApp) {
	spec := brokerapi.DeprovisionServiceSpec{}

	cfClient, err := client.GetClient()
	if err != nil {
		pipeline <- result.NewDeleteApp().WithError(err)
		return
	}

	app := node.App

	_, _, err = cfClient.UpdateApplicationStop(app.GUID)
	if err != nil {
		pipeline <- result.NewDeleteApp().WithError(err)
		return
	}

	_, _, err = cfClient.DeleteRoute(node.Route.GUID)
	if err != nil {
		pipeline <- result.NewDeleteApp().WithError(err)
	}

	_, _, err = cfClient.DeleteApplication(app.GUID)
	if err != nil {
		pipeline <- result.NewDeleteApp().WithError(err)
		return
	}

	pipeline <- result.NewDeleteApp().WithNode(node).WithSpec(spec)
}
