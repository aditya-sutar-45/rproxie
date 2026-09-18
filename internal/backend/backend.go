// Package backend
package backend

import "net/url"

type Backend struct {
	ID   string
	Addr string
	URL  *url.URL
}

func New(addr string, id string) (*Backend, error) {
	backendURL, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	return &Backend{
		ID:   id,
		Addr: addr,
		URL:  backendURL,
	}, nil
}
