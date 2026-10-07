// Package backend
package backend

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/aditya-sutar-45/rproxie/internal/utils"
)

type Backend struct {
	ID                string
	Addr              string
	URL               *url.URL
	health            bool
	mu                sync.RWMutex
	logger            *slog.Logger
	healthCheckClient *http.Client
}

func New(addr string, id string, logger *slog.Logger, healthCheckTimeout time.Duration) (*Backend, error) {
	backendURL, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	backend := &Backend{
		ID:     id,
		Addr:   addr,
		URL:    backendURL,
		logger: logger,
		healthCheckClient: &http.Client{
			Timeout: healthCheckTimeout,
		},
	}

	healthStatus := backend.CheckHealth()
	backend.SetHealth(healthStatus)
	if !healthStatus {
		logger.Warn(
			"backend is unhealthy",
			"id", id,
			"addr", addr,
		)
	}

	return backend, nil
}

func (b *Backend) CheckHealth() bool {
	url := fmt.Sprintf("%s/health", b.URL.String())
	resp, err := b.healthCheckClient.Get(url)
	if err != nil {
		return false
	}
	defer utils.CloseResponseBody(resp, b.logger)

	return resp.StatusCode == http.StatusOK
}

func (b *Backend) SetHealth(health bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.health = health
}

func (b *Backend) GetHealth() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.health
}
