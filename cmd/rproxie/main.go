package main

import (
	"log"
	"time"

	"github.com/aditya-sutar-45/rproxie/internal/config"
	"github.com/aditya-sutar-45/rproxie/internal/logger"
	"github.com/aditya-sutar-45/rproxie/internal/rproxy"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	appLogger := logger.New()

	cfg, err := config.Load()
	if err != nil {
		appLogger.Error("could not load config", "error", err)
	}

	rproxy, err := rproxy.New(
		cfg.Port,
		cfg.Backends,
		appLogger,
		time.Second*5,
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
