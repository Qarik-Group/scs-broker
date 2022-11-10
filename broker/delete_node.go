package broker

import (
	brokerapi "github.com/pivotal-cf/brokerapi/v7/domain"
	"github.com/starkandwayne/scs-broker/broker/result"
)

func (broker *SCSBroker) deleteNode(node *deletable, pipeline chan<- *result.DeleteApp) {
	spec := brokerapi.DeprovisionServiceSpec{}

	cfClient, err := broker.GetClient()
	if err != nil {
		pipeline <- result.NewDeleteApp().WithError(err)
		return
	}

	app := node.App

	routes, _, err := cfClient.GetApplicationRoutes(app.GUID)
	if err != nil {
		pipeline <- result.NewDeleteApp().WithError(err)
		return
	}

	_, _, err = cfClient.UpdateApplicationStop(app.GUID)
	if err != nil {
		pipeline <- result.NewDeleteApp().WithError(err)
		return
	}

	for route := range routes {
		_, _, err := cfClient.DeleteRoute(routes[route].GUID)
		if err != nil {
			pipeline <- result.NewDeleteApp().WithError(err)
			return
		}
	}

	_, _, err = cfClient.DeleteApplication(app.GUID)
	if err != nil {
		pipeline <- result.NewDeleteApp().WithError(err)
		return
	}

	pipeline <- result.NewDeleteApp().WithApp(app).WithSpec(spec)
}
