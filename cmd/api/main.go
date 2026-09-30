package main

import (
	"log"
	"net/http"

	"courtsite-url-shortener/internal/shortener"
)

func main() {
	// Initialize dependencies (DI)
	store := shortener.NewMemoryStore()
	svc := shortener.NewService(store)
	h := shortener.NewHandler(svc)

	// Use Go 1.22+ enhanced ServeMux routing
	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", h.HandleShortener)
	mux.HandleFunc("GET /analytics", h.HandleAnalytics)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
