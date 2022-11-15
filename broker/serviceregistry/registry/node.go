package registry

import "code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"

type Node struct {
	App   ccv3.Application
	Route ccv3.Route
}
