package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/okoye-dev/flux-server/internal/config"
	"github.com/okoye-dev/flux-server/internal/services"
	"github.com/okoye-dev/flux-server/internal/transport/rest"
)

func main() {
	// Load environment variables from .env file (if it exists)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()
	
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Initialize WhatsApp bot if enabled
	if cfg.WhatsApp.Enabled {
		if cfg.WhatsApp.InstanceID == "" || cfg.WhatsApp.Token == "" {
			log.Fatal("WhatsApp bot is enabled but missing required credentials (WHATSAPP_INSTANCE_ID or WHATSAPP_TOKEN)")
		}
		
		log.Printf("Initializing WhatsApp bot with Instance ID: %s", cfg.WhatsApp.InstanceID)
		whatsappBot := services.NewWhatsAppBot(cfg.WhatsApp.InstanceID, cfg.WhatsApp.Token)
		go whatsappBot.Start() // Start bot in a goroutine
		log.Println("WhatsApp bot started successfully and listening for messages...")
	} else {
		log.Println("WhatsApp bot is disabled")
	}

	// Create server with security middleware
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: rest.NewSecureRouter(),
	}

	log.Printf("Starting Flux server on :%s", cfg.Server.Port)
	log.Printf("Environment: %s", cfg.Server.Environment)
	log.Printf("Health check available at: http://localhost:%s/health", cfg.Server.Port)
	log.Printf("Authentication endpoints (username/password only):")
	log.Printf("  - POST /auth/signup")
	log.Printf("  - POST /auth/signin")
	log.Printf("Protected endpoints:")
	log.Printf("  - GET /profile (requires authentication)")
	log.Printf("  - GET /protected (requires authentication)")
	
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}