package utilities

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	brokerapi "github.com/pivotal-cf/brokerapi/domain"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	passwordLength = 30
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// Generate a random password.
func GenClientPassword() string {
	b := make([]byte, passwordLength)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// Given a broker request details object, determine
// which of the supported services the request is
// about.
func GetKind(details interface{}) (string, error) {
	// fun fact: this whole function is gross.
	if d, ok := details.(brokerapi.ProvisionDetails); ok {
		return d.ServiceID, nil
	}

	if d, ok := details.(brokerapi.DeprovisionDetails); ok {
		return d.ServiceID, nil
	}

	if d, ok := details.(brokerapi.BindDetails); ok {
		return d.ServiceID, nil
	}

	if d, ok := details.(brokerapi.UnbindDetails); ok {
		return d.ServiceID, nil
	}

	if d, ok := details.(brokerapi.PollDetails); ok {
		return d.ServiceID, nil
	}

	if d, ok := details.(brokerapi.UpdateDetails); ok {
		return d.ServiceID, nil
	}

	return "", errors.New("service kind not recognized")
}

// Generate a UAA client ID binding name based on the kind of
// service in question and the binding's  ID.
func MakeClientIdForBinding(serviceId string, bindingId string) string {
	return serviceId + "-binding-" + strings.Replace(bindingId, serviceId+"-", "", 1)
}

// Generate an app name based on the kind of service in question
// and a service instance ID.
func MakeAppName(serviceId string, instanceId string) string {
	return serviceId + "-" + instanceId
}
