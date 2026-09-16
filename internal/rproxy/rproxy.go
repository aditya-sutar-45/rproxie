// Package rproxy
package rproxy

import (
	"io"
	"log"
	"net/http"

	"github.com/aditya-sutar-45/rproxie/internal/backend"
	"github.com/aditya-sutar-45/rproxie/internal/utils"
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

func (rp *ReverseProxy) handler(w http.ResponseWriter, r *http.Request) {
	request, err := rp.newRequest(r, rp.backends[0])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadGateway, "backend unavailable")
		return
	}

	resp, err := rp.client.Do(request)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadGateway, "backend unavailable")
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	rp.writeResponse(w, resp)
}

func (rp *ReverseProxy) newRequest(r *http.Request, backend *backend.Backend) (*http.Request, error) {
	targetURL := *backend.URL
	targetURL.Path = r.URL.Path
	targetURL.RawQuery = r.URL.RawQuery

	request, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
	if err != nil {
		return nil, err
	}

	request.Header = r.Header.Clone()

	return request, nil
}

func (rp *ReverseProxy) writeResponse(w http.ResponseWriter, resp *http.Response) {
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, err := io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("ERROR streaming response to client: %v", err)
	}
}
