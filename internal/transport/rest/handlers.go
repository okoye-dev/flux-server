package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/okoye-dev/flux-server/internal/middleware"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}

// HealthHandler handles health check requests
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Service:   "flux-server",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// ProfileHandler handles user profile requests (protected route)
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	userEmail, _ := middleware.GetUserEmail(r)

	response := map[string]interface{}{
		"user_id":    userID,
		"user_email": userEmail,
		"message":    "Welcome to your profile",
		"timestamp":  time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// ProtectedDataHandler handles protected data requests
func ProtectedDataHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":      "This is protected data",
		"user_id":   userID,
		"timestamp": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// NewRouter creates and returns a new HTTP router with all routes
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	
	// Public endpoints
	mux.HandleFunc("/health", HealthHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Welcome to the Flux Server",
			"version": "1.0.0",
		})
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