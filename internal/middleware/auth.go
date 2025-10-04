package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// UserContextKey is the key used to store user information in the request context
type UserContextKey string

const (
	UserIDKey    UserContextKey = "user_id"
	UserEmailKey UserContextKey = "user_email"
)

// SupabaseClaims represents the JWT claims from Supabase
type SupabaseClaims struct {
	Aud   string `json:"aud"`
	Exp   int64  `json:"exp"`
	Iat   int64  `json:"iat"`
	Iss   string `json:"iss"`
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates Supabase JWT tokens and adds user info to context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Get Supabase configuration from environment
		supabaseURL := os.Getenv("SUPABASE_URL")
		supabaseAnonKey := os.Getenv("SUPABASE_ANON_KEY")

		if supabaseURL == "" || supabaseAnonKey == "" {
			http.Error(w, "Supabase configuration missing", http.StatusInternalServerError)
			return
		}

		// Parse and validate the JWT token
		claims, err := validateSupabaseToken(token, supabaseURL, supabaseAnonKey)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Extract username from email (remove @fluxapp.com)
		username := claims.Email
		if strings.HasSuffix(username, "@fluxapp.com") {
			username = strings.TrimSuffix(username, "@fluxapp.com")
		}

		// Add user information to the request context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.Sub)
		ctx = context.WithValue(ctx, UserEmailKey, username) // Using username for compatibility

		// Continue with the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// validateSupabaseToken validates a Supabase JWT token
func validateSupabaseToken(tokenString, supabaseURL, supabaseAnonKey string) (*SupabaseClaims, error) {
	// Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET not configured")
	}

	// Parse and validate the JWT token
	token, err := jwt.ParseWithClaims(tokenString, &SupabaseClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		
		// Use the JWT secret for verification
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*SupabaseClaims); ok && token.Valid {
		// Additional validation for Supabase tokens
		if claims.Iss != supabaseURL+"/auth/v1" {
			return nil, fmt.Errorf("invalid issuer")
		}
		
		// Validate audience
		if claims.Aud != "authenticated" {
			return nil, fmt.Errorf("invalid audience")
		}
		
		// Validate role
		if claims.Role != "authenticated" {
			return nil, fmt.Errorf("invalid role")
		}
		
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GetUserID extracts the user ID from the request context
func GetUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	return userID, ok
}

// GetUserEmail extracts the user email from the request context
func GetUserEmail(r *http.Request) (string, bool) {
	userEmail, ok := r.Context().Value(UserEmailKey).(string)
	return userEmail, ok
}

// OptionalAuthMiddleware is similar to AuthMiddleware but doesn't require authentication
// It adds user info to context if a valid token is provided, but doesn't fail if no token is provided
func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		
		// If no auth header, continue without user context
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			next.ServeHTTP(w, r)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		supabaseURL := os.Getenv("SUPABASE_URL")
		supabaseAnonKey := os.Getenv("SUPABASE_ANON_KEY")

		if supabaseURL == "" || supabaseAnonKey == "" {
			next.ServeHTTP(w, r)
			return
		}

		claims, err := validateSupabaseToken(token, supabaseURL, supabaseAnonKey)
		
		// If token is invalid, continue without user context
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// Extract username from email (remove @fluxapp.com)
		username := claims.Email
		if strings.HasSuffix(username, "@fluxapp.com") {
			username = strings.TrimSuffix(username, "@fluxapp.com")
		}

		// Add user information to context if token is valid
		ctx := context.WithValue(r.Context(), UserIDKey, claims.Sub)
		ctx = context.WithValue(ctx, UserEmailKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}