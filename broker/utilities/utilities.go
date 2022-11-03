package utilities

import (
	"encoding/json"
	"errors"
	"fmt"
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

func ExtractRegistryParams(details string) (*RegistryParams, error) {
	// just in case we got an empty payload
	if len(details) == 0 {
		details = `{"count" : 1}`
	}

	rp := NewRegistryParams()

	if err := json.Unmarshal([]byte(details), rp); err != nil {
		return nil, err
	}

	return rp, nil
}

func NewRegistryParams() *RegistryParams {
	return &RegistryParams{RawCount: 1}
}

type RegistryParams struct {
	RawCount                 int    `json:"count"`
	ApplicationSecurityGroup string `json:"application_security_group"`
}

func (rp *RegistryParams) Count() (int, error) {
	var err error = nil

	if rp.RawCount < 1 {
		err = fmt.Errorf("invalid node count: %d", rp.RawCount)
	}

	return rp.RawCount, err
}
