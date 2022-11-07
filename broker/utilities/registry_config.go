package utilities

import (
	"encoding/json"
	"strings"
)

func NewRegistryConfig() *RegistryConfig {
	rc := &RegistryConfig{}
	rc.Standalone()

	return rc
}

type RegistryConfig struct {
	Mode  string
	Peers []string
}

func (rc *RegistryConfig) AddPeer(peer string) {
	rc.Peers = append(rc.Peers, peer)
}

func (rc *RegistryConfig) Standalone() {
	rc.Mode = "standalone"
}

func (rc *RegistryConfig) Clustered() {
	rc.Mode = "clustered"
}

func (rc *RegistryConfig) String() string {
	return string(rc.Bytes())
}

func (rc *RegistryConfig) Bytes() []byte {
	client := make(map[string]interface{})
	m := rc.Mode == "clustered"

	client["registerWithEureka"] = m
	client["fetchRegistry"] = m

	if len(rc.Peers) > 0 {
		serviceUrl := make(map[string]interface{})
		defaultZone := strings.Join(rc.Peers, ",")
		serviceUrl["defaultZone"] = defaultZone
		client["serviceUrl"] = serviceUrl
	}

	eureka := make(map[string]interface{})
	eureka["client"] = client

	data := make(map[string]interface{})
	data["eureka"] = eureka

	output, err := json.Marshal(data)
	if err != nil {
		return []byte("{}")
	}

	return output

}
