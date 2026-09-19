package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type response struct {
	Message   string `json:"message"`
	BackendID string `json:"backend_id"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Host      string `json:"host"`
	Query     string `json:"query"`
}

type healthResponse struct {
	Status    string `json:"status"`
	BackendID string `json:"backend_id"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	backendID := os.Getenv("BACKEND_ID")
	if backendID == "" {
		backendID = "unknown"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(healthResponse{
			Status:    "ok",
			BackendID: backendID,
		})
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response{
			Message:   "response from test backend",
			BackendID: backendID,
			Method:    r.Method,
			Path:      r.URL.Path,
			Host:      r.Host,
			Query:     r.URL.RawQuery,
		})
	})

	addr := ":" + port
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	shutdownSignal, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("test backend %s listening on %s", backendID, addr)
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
		return
	case <-shutdownSignal.Done():
		log.Printf("test backend %s shutting down", backendID)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("test backend %s shutdown error: %v", backendID, err)
	}
}
