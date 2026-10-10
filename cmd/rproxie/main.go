package main

import (
	"time"

	"github.com/aditya-sutar-45/rproxie/internal/config"
	"github.com/aditya-sutar-45/rproxie/internal/logger"
	"github.com/aditya-sutar-45/rproxie/internal/rproxy"
)

func main() {
	appLogger := logger.New()

	cfg, err := config.LoadYAML("./config_sample.yaml")
	if err != nil {
		appLogger.Error("could not load config", "error", err)
		return
	}

	rproxy, err := rproxy.New(
		cfg.Server.Port,
		cfg.Backends,
		time.Second*5,
		cfg.RateLimiter.BucketCapacity,
		cfg.RateLimiter.RefillRatePerSecond,
		appLogger,
	)
	if err != nil {
		appLogger.Error("creating a reverse proxy", "error", err)
		return
	}

	if err := rproxy.Start(); err != nil {
		appLogger.Error("starting a reverse proxy server", "error", err)
		return
	}
}
