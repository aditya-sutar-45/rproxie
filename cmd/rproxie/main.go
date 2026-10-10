package main

import (
	"fmt"
	"time"

	"github.com/aditya-sutar-45/rproxie/internal/config"
	"github.com/aditya-sutar-45/rproxie/internal/logger"
	"github.com/aditya-sutar-45/rproxie/internal/rproxy"
)

func main() {
	cfg, err := config.Load("./config_sample.yaml")
	if err != nil {
		fmt.Printf("could not load config: %v", err)
		return
	}
	if err := cfg.Validate(); err != nil {
		fmt.Printf("validating config file: %v", err)
		return
	}

	appLogger := logger.New(cfg.Logging.Level)
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
