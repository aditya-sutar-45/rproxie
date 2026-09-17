package main

import (
	"log"

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

	backends := []string{"http://localhost:9000"}

	rproxy, err := rproxy.New(
		8000,
		backends,
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
