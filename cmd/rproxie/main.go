package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/aditya-sutar-45/rproxie/internal/config"
	"github.com/aditya-sutar-45/rproxie/internal/logger"
	"github.com/aditya-sutar-45/rproxie/internal/rproxy"
)

func main() {
	configPath := flag.String("config", "./config_sample.yaml", "path to the config file")
	flag.Parse()

	fmt.Printf("config path: %s\n", *configPath)

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Printf("could not load config: %v", err)
		return
	}
	if err := cfg.Validate(); err != nil {
		fmt.Printf("validating config file: %v", err)
		return
	}

	appLogger := logger.New(cfg.Logging.Level)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	rproxy, err := rproxy.New(
		cfg.Server.Port,
		cfg.Backends,
		time.Second*5,
		cfg.RateLimiter.Enabled,
		cfg.RateLimiter.BucketCapacity,
		cfg.RateLimiter.RefillRatePerSecond,
		appLogger,
	)
	if err != nil {
		appLogger.Error("creating a reverse proxy", "error", err)
		return
	}

	if err := rproxy.Start(ctx); err != nil {
		appLogger.Error("starting a reverse proxy server", "error", err)
		return
	}
}
