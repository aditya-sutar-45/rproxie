// Package rproxy
package rproxy

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/aditya-sutar-45/rproxie/internal/backend"
	"github.com/aditya-sutar-45/rproxie/internal/config"
	"github.com/aditya-sutar-45/rproxie/internal/healthchecker"
	"github.com/aditya-sutar-45/rproxie/internal/loadbalancer"
	"github.com/aditya-sutar-45/rproxie/internal/ratelimiter"
)

type ReverseProxy struct {
	port          int
	backends      []*backend.Backend
	client        *http.Client
	logger        *slog.Logger
	loadBalancer  loadbalancer.LoadBalancer
	healthChecker *healthchecker.HealthChecker
	rateLimiter   *ratelimiter.RateLimiter
}

func New(
	cfg *config.RProxieConfig,
	tickerDuration time.Duration,
	logger *slog.Logger,
) (*ReverseProxy, error) {
	backends := []*backend.Backend{}
	for _, b := range cfg.Backends {
		backend, err := backend.New(
			b.URL,
			b.ID,
			logger,
			2*time.Second,
		)
		if err != nil {
			return nil, err
		}
		backends = append(backends, backend)
	}

	limiter := ratelimiter.New(
		cfg.RateLimiter.BucketCapacity,
		float64(cfg.RateLimiter.RefillRatePerSecond),
		cfg.RateLimiter.Enabled,
	)

	return &ReverseProxy{
		port:          cfg.Server.Port,
		backends:      backends,
		client:        &http.Client{},
		logger:        logger,
		loadBalancer:  loadbalancer.New(len(backends)),
		healthChecker: healthchecker.New(tickerDuration, backends, logger),
		rateLimiter:   limiter,
	}, nil
}

func (rp *ReverseProxy) Start(ctx context.Context) error {
	portString := fmt.Sprintf(":%d", rp.port)
	listener, err := net.Listen("tcp", portString)
	if err != nil {
		rp.logger.Error("error binding port", "port", rp.port, "error", err)
		return err
	}
	defer func() {
		_ = listener.Close()
	}()

	err = rp.serve(ctx, listener)
	if err != nil {
		rp.logger.Error("error serving rproxie", "error", err)
		return err
	}

	return nil
}

func (rp *ReverseProxy) serve(ctx context.Context, listener net.Listener) error {
	healthCheckCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := http.NewServeMux()
	mux.HandleFunc("/", rp.handler)

	server := http.Server{
		Addr:    listener.Addr().String(),
		Handler: mux,
	}

	rp.logger.Info("server starting", "port", rp.port)

	go rp.healthChecker.Start(healthCheckCtx)
	errChannel := make(chan error, 1)

	go func() {
		errChannel <- server.Serve(listener)
	}()

	var serveErr error
	select {
	case <-ctx.Done():
	case serveErr = <-errChannel:
		if errors.Is(serveErr, http.ErrServerClosed) {
			rp.logger.Info("server closed")
			return nil
		}
	}

	rp.logger.Info("server stopped", "port", rp.port)

	shutdown, cancel := context.WithTimeout(
		context.Background(),
		20*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdown); err != nil {
		rp.logger.Error("server shutdown failed", "error", err)
		return err
	}
	rp.logger.Info("server shutdown complete", "port", rp.port)
	return serveErr
}
