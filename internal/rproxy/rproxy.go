// Package rproxy
package rproxy

import (
	"log"
	"net/http"
	"net/url"
)

type ReverseProxy struct {
	portString string
	backendAPI *url.URL
	client     *http.Client
}

func New(portString string, backendAPI string) (*ReverseProxy, error) {
	backendURL, err := url.Parse(backendAPI)
	if err != nil {
		return nil, err
	}
	return &ReverseProxy{
		portString: portString,
		backendAPI: backendURL,
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
