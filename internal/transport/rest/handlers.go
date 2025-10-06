package rest

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/okoye-dev/flux-server/internal/middleware"
	"github.com/okoye-dev/flux-server/internal/models"
	"github.com/okoye-dev/flux-server/internal/services"
)

// HealthHandler handles health check requests
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Service:   "flux-server",
	}

	WriteHealthResponse(w, response)
}


// ProfileHandler handles user profile requests (protected route)
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		WriteInternalServerError(w, MsgUserIDNotFound, "")
		return
	}

	userEmail, _ := middleware.GetUserEmail(r)

	// Fetch profile data from user_profiles table
	profileService, err := services.NewProfileService()
	if err != nil {
		WriteInternalServerError(w, "Failed to create profile service", err.Error())
		return
	}

	profile, err := profileService.GetUserProfile(userID)
	if err != nil {
		// If profile not found, create a placeholder
		authUserID, _ := uuid.Parse(userID)
		profile = &models.UserProfile{
			AuthUserID:  &authUserID,
			DisplayName: &userEmail, // Default to username
			Metadata:    map[string]any{},
		}
	}

	response := ProfileResponse{
		UserID:    userID,
		UserEmail: userEmail,
		Message:   MsgWelcomeToProfile,
		Timestamp: time.Now(),
		Profile:   *profile,
	}

	WriteProfileResponse(w, response)
}

// ProtectedDataHandler handles protected data requests
func ProtectedDataHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		WriteInternalServerError(w, MsgUserIDNotFound, "")
		return
	}

	response := ProtectedDataResponse{
		Data:      MsgProtectedData,
		UserID:    userID,
		Timestamp: time.Now(),
	}

	WriteProtectedDataResponse(w, response)
}

// WhatsAppWebhookHandler handles incoming WhatsApp webhooks
func WhatsAppWebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading webhook body: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Parse the webhook data
	var webhookData map[string]interface{}
	if err := json.Unmarshal(body, &webhookData); err != nil {
		log.Printf("Error parsing webhook JSON: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Printf("WhatsApp webhook received: %+v", webhookData)

	// Process the webhook (you can add your bot logic here)
	// For now, just acknowledge receipt
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// NewRouter creates and returns a new HTTP router with all routes
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	
	// Public endpoints
	mux.HandleFunc("/health", HealthHandler)
	mux.HandleFunc("/webhook/whatsapp", WhatsAppWebhookHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := RootResponse{
			Message: MsgWelcomeToFlux,
			Version: "1.0.0",
		}
		WriteRootResponse(w, response)
	})
	
	// Authentication endpoints (username/password only)
	mux.HandleFunc("/auth/signup", SignupHandler)
	mux.HandleFunc("/auth/signin", SigninHandler)
	
	// Protected endpoints (require authentication)
	mux.Handle("/profile", middleware.AuthMiddleware(http.HandlerFunc(ProfileHandler)))
	mux.Handle("/protected", middleware.AuthMiddleware(http.HandlerFunc(ProtectedDataHandler)))
	
	return mux
}

// NewSecureRouter creates a router with security middleware applied
func NewSecureRouter() http.Handler {
	mux := NewRouter()
	
	// Apply security middleware in order
	handler := middleware.SecurityHeadersMiddleware(mux)
	handler = middleware.CORSMiddleware([]string{"http://localhost:3000", "http://localhost:3002", "http://localhost:8080"})(handler)
	handler = middleware.RateLimitMiddleware(60)(handler) // 60 requests per minute
	
	return handler
}