package utilities

func NewRegistryConfig() *RegistryConfig {
	rc := &RegistryConfig{}
	rc.Standalone()

	return rc
}

type RegistryPeer struct {
	Index  int    `json:"index"`
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
}

type RegistryConfig struct {
	Mode  string
	Peers []*RegistryPeer
}

func (rc *RegistryConfig) AddPeer(idx int, scheme string, host string, port int) {
	rc.Peers = append(rc.Peers, &RegistryPeer{Index: idx, Scheme: scheme, Host: host, Port: port})
}

func (rc *RegistryConfig) Standalone() {
	rc.Mode = "standalone"
}

func (rc *RegistryConfig) Clustered() {
	rc.Mode = "clustered"
}

//func (rc *RegistryConfig) String() string {
//return string(rc.Bytes())
//}

//func (rc *RegistryConfig) Bytes() []byte {
//client := make(map[string]interface{})
//m := rc.Mode == "clustered"

//client["registerWithEureka"] = m
//client["fetchRegistry"] = m

//if len(rc.Peers) > 0 {
//serviceUrl := make(map[string]interface{})
//defaultZone := strings.Join(rc.Peers, ",")
//serviceUrl["defaultZone"] = defaultZone
//client["serviceUrl"] = serviceUrl
//}

//eureka := make(map[string]interface{})
//eureka["client"] = client

//data := make(map[string]interface{})
//data["eureka"] = eureka

//output, err := json.Marshal(data)
//if err != nil {
//return []byte("{}")
//}

//return output

//}
