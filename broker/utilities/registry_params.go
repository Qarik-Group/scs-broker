package utilities

import (
	"encoding/json"
	"fmt"
)

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
