package result

import (
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	brokerapi "github.com/pivotal-cf/brokerapi/v7/domain"
)

type DeleteApp struct {
	App   ccv3.Application
	Spec  brokerapi.DeprovisionServiceSpec
	Error error
}

func NewDeleteApp() *DeleteApp {
	return &DeleteApp{
		App:   ccv3.Application{},
		Spec:  brokerapi.DeprovisionServiceSpec{},
		Error: nil,
	}
}

func (result *DeleteApp) WithApp(app ccv3.Application) *DeleteApp {
	result.App = app

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
