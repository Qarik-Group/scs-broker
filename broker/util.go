package broker

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
func genClientPassword() string {
	b := make([]byte, passwordLength)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// Given a broker request details object, determine
// which of the supported services the request is
// about.
func getKind(details interface{}) (string, error) {
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
func makeClientIdForBinding(kind string, bindingId string) string {
	return kind + "-binding-" + strings.Replace(bindingId, kind+"-", "", 1)
}

// Generate an app name based on the kind of service in question
// and a service instance ID.
func makeAppName(kind string, instanceId string) string {
	return kind + "-" + instanceId
}

func extractRegistryParams(details string) (map[string]interface{}, error) {
	// decode the raw params
	decoded := make(map[string]interface{})
	if err := json.Unmarshal([]byte(details), &decoded); err != nil {
		return nil, err
	}

	// get the registry-specific params that affect broker operations
	rp := registryParams{}

	rp.Merge("count", decoded)
	rp.Merge("application-security-groups", decoded)
	for key, _ := range rp {
		fmt.Println(key)
	}

	return rp, nil
}

type registryParams map[string]interface{}

func (rp registryParams) Merge(key string, other map[string]interface{}) {
	if value, found := other[key]; found {
		rp[key] = value
	}
}
