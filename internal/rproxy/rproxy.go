// Package rproxy
package rproxy

import (
	"log"
	"net/http"
)

type ReverseProxy struct {
	portString string
	backendAPI string
}

func New(portString string, backendAPI string) *ReverseProxy {
	return &ReverseProxy{
		portString: portString,
		backendAPI: backendAPI,
	}
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
