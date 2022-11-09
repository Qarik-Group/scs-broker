package hype

type Header struct {
	Name  string
	Value string
}

func NewHeader(name string, value string) *Header {
	return &Header{Name: name, Value: value}
}
