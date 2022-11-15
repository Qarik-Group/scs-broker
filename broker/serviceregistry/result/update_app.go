package result

import (
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
)

type UpdateApp struct {
	App   ccv3.Application
	Error error
}

func NewUpdateApp() *UpdateApp {
	return &UpdateApp{
		App:   ccv3.Application{},
		Error: nil,
	}
}

func (result *UpdateApp) WithApp(app ccv3.Application) *UpdateApp {
	result.App = app

	return result
}

func (result *UpdateApp) WithError(err error) *UpdateApp {
	result.Error = err

	return result
}

func (result *UpdateApp) Failure() error {
	return result.Failure()
}
