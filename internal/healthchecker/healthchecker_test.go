package healthchecker

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aditya-sutar-45/rproxie/internal/backend"
	"github.com/aditya-sutar-45/rproxie/internal/logger"
)

func TestHealthChecker(t *testing.T) {
	logger := logger.New("INFO")

	server := newTestServer(t)
	defer server.Close()

	b1 := newBackend(t, server, logger, 3*time.Second)
	b2 := newBackend(t, server, logger, 3*time.Second)
	b3 := newBackend(t, server, logger, 3*time.Second)

	backends := []*backend.Backend{b1, b2, b3}

	start := time.Now()
	h := New(1*time.Second, backends, logger)

	logger.Info("Starting health check")
	h.startHealthChecks()

	timeTaken := time.Since(start)
	logger.Info("Health check completed", "timeTaken", timeTaken)
}

func newTestServer(t *testing.T) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			t.Errorf("expected /health, got %s", r.URL.Path)
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))

	return server
}

func newBackend(t *testing.T, server *httptest.Server, logger *slog.Logger, timeout time.Duration) *backend.Backend {
	b, err := backend.New(server.URL, "1", logger, timeout)
	if err != nil {
		t.Fatal("Failed to create backend 1:", err)
	}

	return b
}
