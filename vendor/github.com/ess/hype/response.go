package hype

import (
	"net/http"
)

type Response struct {
	actual *http.Response
	data   []byte
	err    error
}

func (response Response) Okay() bool {
	if response.err == nil {
		return true
	}
	return false
}

func (response Response) Data() []byte {
	return response.data
}

func (response Response) Error() error {
	return response.err
}

func (response Response) Header(name string) string {
	return response.actual.Header.Get(name)
}
