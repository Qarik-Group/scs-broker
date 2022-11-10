package utilities

import (
	"encoding/json"
	"fmt"
	"strings"
)

func NewRegistryConfig() *RegistryConfig {
	rc := &RegistryConfig{}
	rc.Standalone()

	return rc
}

type RegistryPeer struct {
	//Index  int    `json:"index"`
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
	//Port   int    `json:"port"`
}

func (peer *RegistryPeer) String() string {
	return fmt.Sprintf("%s://%s", peer.Scheme, peer.Host)
}

type RegistryConfig struct {
	Mode  string
	Peers []*RegistryPeer
}

func (rc *RegistryConfig) AddPeer(scheme string, host string) {
	rc.Peers = append(rc.Peers, &RegistryPeer{Scheme: scheme, Host: host})
}

func (rc *RegistryConfig) Standalone() {
	rc.Mode = "standalone"
}

func (rc *RegistryConfig) Clustered() {
	rc.Mode = "clustered"
}

func (rc *RegistryConfig) ForNode(node string) string {
	client := make(map[string]interface{})
	m := rc.Mode == "clustered"

	client["registerWithEureka"] = m
	client["fetchRegistry"] = m

	if m {
		peers := make([]string, 0)
		for _, peer := range rc.Peers {
			if peer.Host != node {
				peers = append(peers, peer.String())
			}
		}

		serviceUrl := make(map[string]interface{})
		defaultZone := strings.Join(peers, ",")
		serviceUrl["defaultZone"] = defaultZone
		client["serviceUrl"] = serviceUrl
	}

	eureka := make(map[string]interface{})
	eureka["client"] = client

	data := make(map[string]interface{})
	data["eureka"] = eureka

	output, err := json.Marshal(data)
	if err != nil {
		return "{}"
	}

	return string(output)

}
