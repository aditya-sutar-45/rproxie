// Package healthchecker
package healthchecker

import (
	"log/slog"
	"time"

	"github.com/aditya-sutar-45/rproxie/internal/backend"
)

type HealthChecker struct {
	healthCheckDuration time.Duration
	backends            []*backend.Backend
	logger              *slog.Logger
}

func New(duration time.Duration, backends []*backend.Backend, logger *slog.Logger) *HealthChecker {
	return &HealthChecker{
		healthCheckDuration: duration,
		backends:            backends,
		logger:              logger,
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

		for _, b := range h.backends {
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
	}
}
