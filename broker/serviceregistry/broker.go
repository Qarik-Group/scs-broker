package serviceregistry

type Broker struct{}

func New() *Broker {
	return &Broker{}
}
