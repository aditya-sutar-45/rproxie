// Package rproxy
package rproxy

import (
	"log"
	"net/http"

	"github.com/aditya-sutar-45/rproxie/internal/backend"
)

type ReverseProxy struct {
	portString string
	backends   []*backend.Backend
	client     *http.Client
}

func New(portString string, backendAddrs []string) (*ReverseProxy, error) {
	backends := []*backend.Backend{}
	for _, b := range backendAddrs {
		backend, err := backend.New(b)
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
