// Package rproxy
package rproxy

import (
	"log"
	"net/http"
	"net/url"
)

type Backend struct {
	Addr string
	URL  *url.URL
}

func NewBackend(addr string) (*Backend, error) {
	backendURL, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	return &Backend{
		Addr: addr,
		URL:  backendURL,
	}, nil
}

type ReverseProxy struct {
	portString string
	backends   []*Backend
	client     *http.Client
}

func New(portString string, backendAddrs []string) (*ReverseProxy, error) {
	backends := []*Backend{}
	for _, b := range backendAddrs {
		backend, err := NewBackend(b)
		if err != nil {
			return nil, err
		}
		backends = append(backends, backend)
	}

	return &ReverseProxy{
		portString: portString,
		backends:   backends,
		client:     &http.Client{},
	}, nil
}

func (rp *ReverseProxy) Start() error {
	http.HandleFunc("/", rp.handler)

	log.Printf("INFO server is starting on port %s\n", rp.portString)

	err := http.ListenAndServe(rp.portString, nil)
	if err != nil {
		return err
	}

	return nil
}
