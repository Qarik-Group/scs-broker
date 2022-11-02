package broker

type ServiceRegistryInstance struct {
	SCSBroker *SCSBroker
}

type scs_instance interface {
	CreateServerInstance(string, string, string, map[string]string) (string, error)
}
