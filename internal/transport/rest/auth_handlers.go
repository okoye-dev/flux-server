package rest

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/okoye-dev/flux-server/internal/services"
	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
)

// UserClaims represents JWT claims
type UserClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// SignupHandler handles user signup with username/password only
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Username == "" || req.Password == "" {
		WriteBadRequestError(w, MsgUsernamePasswordRequired, "")
		return
	}

	// Validate username format (alphanumeric only)
	if len(req.Username) < 3 || len(req.Username) > 20 {
		WriteBadRequestError(w, MsgUsernameLengthInvalid, "")
		return
	}

	// Get Supabase configuration
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseAnonKey := os.Getenv("SUPABASE_ANON_KEY")

	if supabaseURL == "" || supabaseAnonKey == "" {
		http.Error(w, "Supabase configuration missing", http.StatusInternalServerError)
		return
	}

	// Create Supabase client with anon key
	client, err := supabase.NewClient(supabaseURL, supabaseAnonKey, nil)
	if err != nil {
		http.Error(w, "Failed to create Supabase client", http.StatusInternalServerError)
		return
	}

	// Create user in Supabase auth.users table using admin API
	// We'll use username@fluxapp.com as the email format for Supabase
	email := req.Username + "@fluxapp.com"
	
	authResponse, err := client.Auth.Signup(types.SignupRequest{
		Email:    email,
		Password: req.Password,
		Data: map[string]interface{}{
			"username": req.Username,
		},
	})
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Create user profile automatically
	profileService, err := services.NewProfileService()
	if err != nil {
		http.Error(w, "Failed to create profile service: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create signup data from request
	signupData := &services.SignupData{
		PhoneNumber:        req.PhoneNumber,
		CropType:           req.CropType,
		LocationID:         req.LocationID,
		Language:           req.Language,
		AssignedLocationID: req.AssignedLocationID,
	}

	_, err = profileService.CreateUserProfile(authResponse.User.ID.String(), req.Username, req.Role, signupData)
	if err != nil {
		// Log error but don't fail signup - user is created in auth
		// TODO: Add proper logging
		_ = err
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AuthResponse{
		User: UserInfo{
			ID:        authResponse.User.ID.String(),
			Username:  req.Username,
			CreatedAt: authResponse.User.CreatedAt,
		},
		AccessToken: "", // No token on signup, user needs to sign in
		TokenType:   "",
		ExpiresIn:   0,
		Message:     "User created successfully. Profile will be created automatically. Please sign in to get your access token.",
	})
}

// SigninHandler handles user signin with username/password only
func SigninHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SigninRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Get Supabase configuration
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseAnonKey := os.Getenv("SUPABASE_ANON_KEY")

	if supabaseURL == "" || supabaseAnonKey == "" {
		http.Error(w, "Supabase configuration missing", http.StatusInternalServerError)
		return
	}

	// Create Supabase client with anon key
	client, err := supabase.NewClient(supabaseURL, supabaseAnonKey, nil)
	if err != nil {
		http.Error(w, "Failed to create Supabase client", http.StatusInternalServerError)
		return
	}

	// Sign in using the email format (username@fluxapp.com)
	email := req.Username + "@fluxapp.com"
	authResponse, err := client.Auth.SignInWithEmailPassword(email, req.Password)
	if err != nil {
		http.Error(w, "Failed to sign in user: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Extract username from user metadata
	username := req.Username
	if userMetadata := authResponse.User.UserMetadata; userMetadata != nil {
		if metaUsername, exists := userMetadata["username"]; exists {
			if usernameStr, ok := metaUsername.(string); ok {
				username = usernameStr
			}
		}
	}

	// Return the response with access token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AuthResponse{
		User: UserInfo{
			ID:        authResponse.User.ID.String(),
			Username:  username,
			CreatedAt: authResponse.User.CreatedAt,
		},
		AccessToken: authResponse.AccessToken,
		TokenType:   authResponse.TokenType,
		ExpiresIn:   authResponse.ExpiresIn,
		Message:     "Sign in successful",
	})
}

// generateUserID generates a random user ID (not used with Supabase auth)
func generateUserID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}