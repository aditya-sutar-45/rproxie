package healthchecker

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aditya-sutar-45/rproxie/internal/backend"
	"github.com/aditya-sutar-45/rproxie/internal/logger"
)

func TestHealthChecker(t *testing.T) {
	logger := logger.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// validation
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	b1, err := backend.New(server.URL, "1", logger)
	if err != nil {
		t.Fatal("Failed to create backend 1:", err)
		return
	}
	b2, err := backend.New(server.URL, "2", logger)
	if err != nil {
		t.Fatal("Failed to create backend 2:", err)
		return
	}
	b3, err := backend.New(server.URL, "3", logger)
	if err != nil {
		t.Fatal("Failed to create backend 3:", err)
		return
	}

	backends := []*backend.Backend{b1, b2, b3}

	start := time.Now()
	h := New(1*time.Second, backends, logger)

	logger.Info("Starting health check")
	h.startHealthChecks()

	timeTaken := time.Since(start)
	logger.Info("Health check completed", "timeTaken", timeTaken)
}
