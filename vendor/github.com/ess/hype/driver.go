package hype

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ess/debuggable"
)

type Driver struct {
	raw     *http.Client
	baseURL url.URL
}

func New(baseURL string) (*Driver, error) {
	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	d := &Driver{
		&http.Client{Timeout: 20 * time.Second},
		*url,
	}

	return d, nil
}

func (driver *Driver) WithoutTLSVerification() *Driver {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Timeout: driver.raw.Timeout, Transport: tr}

	return &Driver{client, driver.baseURL}
}

func (driver *Driver) Get(path string, params Params) *Request {
	return driver.newRequest("GET", path, saneParams(params).ToValues(), nil)
}

func (driver *Driver) Post(path string, params Params, data []byte) *Request {
	return driver.newRequest("POST", path, saneParams(params).ToValues(), data)
}

func (driver *Driver) Put(path string, params Params, data []byte) *Request {
	return driver.newRequest("PUT", path, saneParams(params).ToValues(), data)
}

func (driver *Driver) Patch(path string, params Params, data []byte) *Request {
	return driver.newRequest("PATCH", path, saneParams(params).ToValues(), data)
}

func (driver *Driver) Delete(path string, params Params) *Request {
	return driver.newRequest("DELETE", path, saneParams(params).ToValues(), nil)
}

func (driver *Driver) newRequest(verb string, path string, params url.Values, data []byte) *Request {
	request, err := http.NewRequest(
		verb,
		driver.constructRequestURL(path, params),
		bytes.NewReader(data),
	)

	if err != nil {
		return &Request{nil, nil, err}
	}

	return &Request{request, driver.raw, nil}
}

func (driver *Driver) constructRequestURL(path string, params url.Values) string {

	pathParts := []string{driver.baseURL.Path, path}

	requestURL := url.URL{
		Scheme:   driver.baseURL.Scheme,
		Host:     driver.baseURL.Host,
		Path:     strings.Join(pathParts, "/"),
		RawQuery: params.Encode(),
	}

	result := requestURL.String()

	if debuggable.Enabled() {
		fmt.Println("[DEBUG] Request URL:", result)
	}

	return result
}

/*
Copyright 2021 Dennis Walters

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
