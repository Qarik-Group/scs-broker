package hype

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Request struct {
	actual *http.Request
	raw    *http.Client
	err    error
}

func (request *Request) WithHeader(header *Header) *Request {
	return request.WithHeaderSet(header)
}

func (request *Request) WithHeaderSet(headers ...*Header) *Request {
	for _, header := range headers {
		request.actual.Header.Add(header.Name, header.Value)
	}

	return request
}

func (request *Request) Response() Response {
	response, err := request.raw.Do(request.actual)
	if err != nil {
		return Response{response, nil, err}
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return Response{response, nil, err}
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return Response{response, nil, fmt.Errorf("response status: %d", response.StatusCode)}
	}

	return Response{response, body, nil}
}
