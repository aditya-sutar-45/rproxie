package main

import (
	"log"

	"github.com/aditya-sutar-45/rproxie/internal/rproxy"
)

func main() {
	rproxy, err := rproxy.New(
		":8080",
		"http://localhost:9000",
	)
	if err != nil {
		log.Fatalf("ERROR creating a reverse proxy: %v", err)
	}

	if err := rproxy.Start(); err != nil {
		log.Fatalf("error starting proxy: %v", err)
	}
}
