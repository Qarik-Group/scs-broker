package hype

import (
	"net/url"
)

type Params map[string][]string

func (params Params) Set(key string, value string) {
	params[key] = []string{value}
}

func (params Params) ToValues() url.Values {
	return url.Values(params)
}

func saneParams(p Params) Params {
	if p == nil {
		return make(Params)
	}

	return p
}
