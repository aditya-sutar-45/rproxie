// Package rproxy
package rproxy

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/aditya-sutar-45/rproxie/internal/backend"
	"github.com/aditya-sutar-45/rproxie/internal/healthchecker"
	"github.com/aditya-sutar-45/rproxie/internal/loadbalancer"
	"github.com/aditya-sutar-45/rproxie/internal/utils"
)

type ReverseProxy struct {
	port          int
	backends      []*backend.Backend
	client        *http.Client
	logger        *slog.Logger
	loadBalancer  loadbalancer.LoadBalancer
	healthChecker *healthchecker.HealthChecker
}

func New(port int, backendAddrs []string, logger *slog.Logger, tickerDuration time.Duration) (*ReverseProxy, error) {
	backends := []*backend.Backend{}
	for i, b := range backendAddrs {
		backend, err := backend.New(b, strconv.Itoa(i), logger)
		if err != nil {
			return nil, err
		}
		backends = append(backends, backend)
	}

	return &ReverseProxy{
		port:          port,
		backends:      backends,
		client:        &http.Client{},
		logger:        logger,
		loadBalancer:  loadbalancer.New(len(backends)),
		healthChecker: healthchecker.New(tickerDuration, backends, logger),
	}, nil
}

func (rp *ReverseProxy) Start() error {
	portString := fmt.Sprintf(":%d", rp.port)
	http.HandleFunc("/", rp.handler)

	rp.logger.Info("server starting", "port", rp.port)

	go rp.healthChecker.Start()

	return http.ListenAndServe(portString, nil)
}

func (rp *ReverseProxy) handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var logErr error

	defer func() {
		if logErr != nil {
			rp.logger.Error(
				"backend request failed",
				"method", r.Method,
				"path", r.URL.Path,
				"error", logErr,
			)
			return
		}

		rp.logger.Info(
			"request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"duration", time.Since(start),
		)
	}()

	backendIndex, err := rp.getNextHealthyBackend()
	if err != nil {
		logErr = err
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal Application Error")
		return
	}

	request, err := rp.newRequest(r, rp.backends[backendIndex])
	if err != nil {
		logErr = err
		utils.RespondWithError(w, http.StatusBadGateway, "backend unavailable")
		return
	}

	resp, err := rp.client.Do(request)
	if err != nil {
		logErr = err
		utils.RespondWithError(w, http.StatusBadGateway, "backend unavailable")
		return
	}
	defer utils.CloseResponseBody(resp, rp.logger)

	logErr = rp.writeResponse(w, resp)
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

func (rp *ReverseProxy) writeResponse(w http.ResponseWriter, resp *http.Response) error {
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	_, err := io.Copy(w, resp.Body)
	return err
}

func (rp *ReverseProxy) getNextHealthyBackend() (int, error) {
	backendIndex := -1
	for range rp.backends {
		index, err := rp.loadBalancer.NextBackendIndex()
		if err != nil {
			return backendIndex, err
		}
		if rp.isBackendHealthy(index) {
			backendIndex = index
			break
		}
	}

	if backendIndex == -1 {
		return backendIndex, fmt.Errorf("all backends are unhealthy, backend count: %d", len(rp.backends))
	}

	return backendIndex, nil
}

func (rp *ReverseProxy) isBackendHealthy(index int) bool {
	return rp.backends[index].GetHealth()
}
