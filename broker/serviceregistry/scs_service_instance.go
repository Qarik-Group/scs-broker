package serviceregistry

type ServiceRegistryInstance struct {
	Broker *Broker
}

type scs_instance interface {
	CreateServerInstance(string, string, string, map[string]string) (string, error)
}
