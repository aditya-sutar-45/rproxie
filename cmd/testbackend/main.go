package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type response struct {
	Message   string `json:"message"`
	BackendID string `json:"backend_id"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Host      string `json:"host"`
	Query     string `json:"query"`
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
	log.Printf("test backend %s listening on %s", backendID, addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
