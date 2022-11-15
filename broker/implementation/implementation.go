package implementation

import (
	"context"

	"github.com/pivotal-cf/brokerapi/v7/domain"
	"github.com/starkandwayne/scs-broker/config"
)

type Implementation interface {
	Provision(context.Context, string, domain.ProvisionDetails, bool) (domain.ProvisionedServiceSpec, error)
	Deprovision(context.Context, string, domain.DeprovisionDetails, bool) (domain.DeprovisionServiceSpec, error)
	Update(context.Context, string, domain.UpdateDetails, bool) (domain.UpdateServiceSpec, error)
	Bind(context.Context, string, string, domain.BindDetails, bool) (domain.Binding, error)
	Unbind(context.Context, string, string, domain.UnbindDetails, bool) (domain.UnbindSpec, error)
}

var implementations map[string]Implementation
var unknown Implementation

func init() {
	implementations = make(map[string]Implementation)
	unknown = &failsafe{}
}

func Register(topic string, imp Implementation) {
	implementations[topic] = imp
}

func For(topic string) Implementation {
	if imp, ok := implementations[topic]; ok {
		return imp
	}

	return unknown
}

func ByServiceID(serviceID string) Implementation {
	for key, service := range config.Parsed.Services {
		if service.ServiceID == serviceID {
			return For(key)
		}
	}

	return unknown
}
