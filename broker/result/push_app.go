package result

import (
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
)

type PushApp struct {
	App   ccv3.Application
	Route ccv3.Route
	Error error
}

func NewPushApp() *PushApp {
	return &PushApp{
		App:   ccv3.Application{},
		Route: ccv3.Route{},
		Error: nil,
	}
}

func (result *PushApp) WithApp(app ccv3.Application) *PushApp {
	result.App = app

	return result
}

func (result *PushApp) WithRoute(rte ccv3.Route) *PushApp {
	result.Route = rte

	return result
}

func (result *PushApp) WithError(err error) *PushApp {
	result.Error = err

	return result
}
