package main

import (
	"log"
	"net/http"

	"github.com/okoye-dev/flux-server/internal/transport/rest"
)

func main() {
	server := &http.Server{
		Addr:    ":8080",
		Handler: rest.NewRouter(),
	}

	log.Println("Starting Flux server on :8080")
	log.Println("Health check available at: http://localhost:8080/health")
	
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}