package rproxy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aditya-sutar-45/rproxie/internal/config"
	"github.com/aditya-sutar-45/rproxie/internal/logger"
	"github.com/aditya-sutar-45/rproxie/internal/utils"
)

func TestReverseProxyIntegration(t *testing.T) {
	//  Create a test backend.
	testBackend := newTestBackend()
	defer testBackend.Close()

	//  Create a ReverseProxy configured to use that backend.
	cfg := &config.RProxieConfig{
		Server: config.ServerConfig{
			Port: 6969,
		},
		Backends: []*config.BackendConfig{
			{
				ID:  "test-backend",
				URL: testBackend.URL,
			},
		},
		RateLimiter: &config.RateLimiterConfig{
			Enabled:             true,
			BucketCapacity:      10,
			RefillRatePerSecond: 2,
		},
		Logging: &config.LoggingConfig{
			Level: "INFO",
		},
	}

	testLogger := logger.New("INFO")
	rp, err := New(cfg, 2*time.Second, testLogger)
	if err != nil {
		t.Errorf("error creating a reverse proxy: %v", err)
		return
	}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	done := make(chan error, 1)

	go func() {
		done <- rp.Start(ctx)
	}()

	//  Send an HTTP request through the proxy.
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	proxyURL := fmt.Sprintf(
		"http://127.0.0.1:%d/inspect?source=test",
		cfg.Server.Port,
	)

	var resp *http.Response
	deadline := time.Now().Add(5 * time.Second)

	for time.Now().Before(deadline) {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			proxyURL,
			nil,
		)
		if err != nil {
			t.Fatal(err)
		}

		resp, err = client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}

		if resp != nil {
			utils.CloseResponseBody(resp, testLogger)
			resp = nil
		}

		time.Sleep(50 * time.Millisecond)
	}

	if resp == nil {
		t.Fatal("proxy did not return a successful response")
	}
	defer utils.CloseResponseBody(resp, testLogger)

	//  Verify the response.
	var got struct {
		Method string `json:"method"`
		Path   string `json:"path"`
		Query  string `json:"query"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decoding proxy response: %v", err)
	}

	if got.Method != http.MethodGet {
		t.Errorf("method = %q, want GET", got.Method)
	}

	if got.Path != "/inspect" {
		t.Errorf("path = %q, want /inspect", got.Path)
	}

	if got.Query != "source=test" {
		t.Errorf("query = %q, want source=test", got.Query)
	}
	//  Cancel the context and wait for shutdown.
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("proxy shutdown failed: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("proxy did not shut down in time")
	}
}

func newTestBackend() *httptest.Server {
	testBackend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_ = json.NewEncoder(w).Encode(map[string]any{
			"method":  r.Method,
			"path":    r.URL.Path,
			"query":   r.URL.RawQuery,
			"host":    r.Host,
			"headers": r.Header,
		})
	}))
	return testBackend
}
