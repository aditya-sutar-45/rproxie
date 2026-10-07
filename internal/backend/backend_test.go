package backend

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aditya-sutar-45/rproxie/internal/logger"
)

func TestCheckHealthTimeout(t *testing.T) {
	logger := logger.New()
	server := newTestServer(t)
	defer server.Close()

	start := time.Now()
	b, err := New(server.URL, "1", logger, 1*time.Second)
	if err != nil {
		t.Fatal("Failed to create backend:", err)
	}
	end := time.Since(start)

	currentHealth := b.GetHealth()

	if currentHealth {
		t.Errorf("expected backend to be unhealthy due to timeout, but got healthy")
	}

	logger.Info("Health check completed", "timeTaken", end, "status", currentHealth)
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
