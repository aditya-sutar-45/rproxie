// Package healthchecker
package healthchecker

import (
	"log/slog"
	"sync"
	"time"

	"github.com/aditya-sutar-45/rproxie/internal/backend"
)

type HealthChecker struct {
	healthCheckDuration time.Duration
	backends            []*backend.Backend
	logger              *slog.Logger
	wg                  sync.WaitGroup
}

func New(duration time.Duration, backends []*backend.Backend, logger *slog.Logger) *HealthChecker {
	return &HealthChecker{
		healthCheckDuration: duration,
		backends:            backends,
		logger:              logger,
		wg:                  sync.WaitGroup{},
	}
}

func (h *HealthChecker) Start() {
	ticker := time.NewTicker(h.healthCheckDuration)
	defer ticker.Stop()

	for range ticker.C {

		h.logger.Info(
			"performing health checks",
			"backendCount", len(h.backends),
		)

		h.startHealthChecks()
	}
}

func (h *HealthChecker) startHealthChecks() {
	for _, b := range h.backends {
		h.wg.Add(1)
		go h.performHealthCheck(b)
	}

	h.wg.Wait()
}

func (h *HealthChecker) performHealthCheck(b *backend.Backend) {
	defer h.wg.Done()

	currHealthStatus := b.GetHealth()
	healthStatus := b.CheckHealth()

	if currHealthStatus != healthStatus {
		b.SetHealth(healthStatus)

		if currHealthStatus && !healthStatus {
			h.logger.Error(
				"backend became unavailable",
				"backendID", b.ID,
				"backendURL", b.URL,
			)
		}
		if !currHealthStatus && healthStatus {
			h.logger.Info(
				"backend recovered",
				"backendID", b.ID,
				"backendURL", b.URL,
			)
		}
	}
}
