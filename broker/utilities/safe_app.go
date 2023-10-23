package utilities

import "code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"

func SafeApp(app ccv3.Application) ccv3.Application {
	return ccv3.Application{
		GUID:                app.GUID,
		StackName:           app.StackName,
		LifecycleBuildpacks: app.LifecycleBuildpacks,
		LifecycleType:       app.LifecycleType,
		Metadata:            app.Metadata,
		Name:                app.Name,
		State:               app.State,
	}
}
