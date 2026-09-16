// Package backend
package backend

import "net/url"

type Backend struct {
	Addr string
	URL  *url.URL
}

func New(addr string) (*Backend, error) {
	backendURL, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	return &Backend{
		Addr: addr,
		URL:  backendURL,
	}, nil
}
