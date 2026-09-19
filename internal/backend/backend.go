// Package backend
package backend

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"sync"

	"github.com/aditya-sutar-45/rproxie/internal/utils"
)

type Backend struct {
	ID     string
	Addr   string
	URL    *url.URL
	health bool
	mu     sync.RWMutex
	logger *slog.Logger
}

func New(addr string, id string, logger *slog.Logger) (*Backend, error) {
	backendURL, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	backend := &Backend{
		ID:     id,
		Addr:   addr,
		URL:    backendURL,
		logger: logger,
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
	resp, err := http.Get(url)
	if err != nil {
		b.logger.Debug("sending request", "url", url, "error", err)
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
