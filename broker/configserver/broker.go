package configserver

type Broker struct{}

func New() *Broker {
	return &Broker{}
}
