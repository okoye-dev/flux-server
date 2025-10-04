package rest

import (
	"time"

	"github.com/okoye-dev/flux-server/internal/models"
)

// Standard HTTP Response Types

// APIResponse represents a standard API response
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     *APIError   `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// APIError represents an API error
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// Authentication Response Types

// AuthResponse represents authentication response
type AuthResponse struct {
	User        UserInfo `json:"user"`
	AccessToken string   `json:"access_token"`
	TokenType   string   `json:"token_type"`
	ExpiresIn   int      `json:"expires_in"`
	Message     string   `json:"message"`
}

// UserInfo represents user information in auth responses
type UserInfo struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// SignupRequest represents signup request
type SignupRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role,omitempty"` // Optional role name: "farmer", "extension_officer", defaults to "farmer"`
	
	// Optional farmer-specific fields
	PhoneNumber string `json:"phone_number,omitempty"`
	CropType    string `json:"crop_type,omitempty"`
	LocationID  int64  `json:"location_id,omitempty"`
	Language    string `json:"language,omitempty"` // Defaults to "en"
	
	// Optional extension officer fields
	AssignedLocationID int64 `json:"assigned_location_id,omitempty"`
}

// SigninRequest represents signin request
type SigninRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Profile Response Types

// ProfileResponse represents profile response
type ProfileResponse struct {
	UserID    string                    `json:"user_id"`
	UserEmail string                    `json:"user_email"`
	Message   string                    `json:"message"`
	Timestamp time.Time                 `json:"timestamp"`
	Profile   models.UserProfile        `json:"profile"`
}

// UserProfileResponse represents a user profile response with role info
type UserProfileResponse struct {
	models.UserProfile
	Role *models.Role `json:"role,omitempty"`
}

// UserProfileWithAuthResponse represents a user profile response with auth info
type UserProfileWithAuthResponse struct {
	models.UserProfile
	AuthUser *models.AuthUser `json:"auth_user,omitempty"`
}

// Protected Data Response Types

// ProtectedDataResponse represents protected data response
type ProtectedDataResponse struct {
	Data      string    `json:"data"`
	UserID    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}

// Model Response Types

// FarmerResponse represents farmer response
type FarmerResponse struct {
	models.Farmer
}

// ExtensionOfficerResponse represents extension officer response
type ExtensionOfficerResponse struct {
	models.ExtensionOfficer
}

// RoleResponse represents role response
type RoleResponse struct {
	models.Role
}

// LocationResponse represents location response
type LocationResponse struct {
	models.Location
}

// CropResponse represents crop response
type CropResponse struct {
	models.Crop
}

// FarmHarvestResponse represents farm harvest response
type FarmHarvestResponse struct {
	models.FarmHarvest
}

// FarmHarvestWithDetailsResponse represents farm harvest with details response
type FarmHarvestWithDetailsResponse struct {
	models.FarmHarvestWithDetails
}

// List Response Types

// FarmersListResponse represents farmers list response
type FarmersListResponse struct {
	Farmers    []models.Farmer `json:"farmers"`
	Pagination Pagination      `json:"pagination"`
}

// ExtensionOfficersListResponse represents extension officers list response
type ExtensionOfficersListResponse struct {
	ExtensionOfficers []models.ExtensionOfficer `json:"extension_officers"`
	Pagination        Pagination                `json:"pagination"`
}

// RolesListResponse represents roles list response
type RolesListResponse struct {
	Roles      []models.Role `json:"roles"`
	Pagination Pagination    `json:"pagination"`
}

// UserProfilesListResponse represents user profiles list response
type UserProfilesListResponse struct {
	UserProfiles []models.UserProfile `json:"user_profiles"`
	Pagination   Pagination           `json:"pagination"`
}

// LocationsListResponse represents locations list response
type LocationsListResponse struct {
	Locations  []models.Location `json:"locations"`
	Pagination Pagination        `json:"pagination"`
}

// CropsListResponse represents crops list response
type CropsListResponse struct {
	Crops      []models.Crop `json:"crops"`
	Pagination Pagination    `json:"pagination"`
}

// FarmHarvestsListResponse represents farm harvests list response
type FarmHarvestsListResponse struct {
	FarmHarvests []models.FarmHarvest `json:"farm_harvests"`
	Pagination   Pagination           `json:"pagination"`
}

// Health Response Types

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}

// Root Response Types

// RootResponse represents root endpoint response
type RootResponse struct {
	Message string `json:"message"`
	Version string `json:"version"`
}

// Error Response Types

// ValidationError represents validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// ValidationErrorResponse represents validation error response
type ValidationErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}

// Common Response Messages
const (
	MsgUserCreatedSuccessfully     = "User created successfully. Profile will be created automatically. Please sign in to get your access token."
	MsgSignInSuccessful           = "Sign in successful"
	MsgWelcomeToProfile           = "Welcome to your profile"
	MsgProtectedData              = "This is protected data"
	MsgServiceHealthy             = "Service is healthy"
	MsgWelcomeToFlux              = "Welcome to the Flux Server"
	MsgInvalidRequest             = "Invalid request"
	MsgUnauthorized               = "Unauthorized"
	MsgForbidden                  = "Forbidden"
	MsgNotFound                   = "Not found"
	MsgInternalServerError        = "Internal server error"
	MsgUsernamePasswordRequired   = "Username and password are required"
	MsgUsernameLengthInvalid      = "Username must be between 3 and 20 characters"
	MsgSupabaseConfigMissing      = "Supabase configuration missing"
	MsgFailedToCreateUser         = "Failed to create user"
	MsgFailedToSignIn             = "Failed to sign in user"
	MsgUserIDNotFound             = "User ID not found in context"
	MsgAuthorizationHeaderRequired = "Authorization header is required"
	MsgInvalidAuthorizationFormat = "Invalid authorization header format"
	MsgInvalidOrExpiredToken      = "Invalid or expired token"
)

// Common Error Codes
const (
	ErrCodeValidation         = "VALIDATION_ERROR"
	ErrCodeUnauthorized       = "UNAUTHORIZED"
	ErrCodeForbidden          = "FORBIDDEN"
	ErrCodeNotFound           = "NOT_FOUND"
	ErrCodeInternalError      = "INTERNAL_ERROR"
	ErrCodeSupabaseError      = "SUPABASE_ERROR"
	ErrCodeAuthError          = "AUTH_ERROR"
	ErrCodeInvalidToken       = "INVALID_TOKEN"
	ErrCodeMissingConfig      = "MISSING_CONFIG"
	ErrCodeUserNotFound       = "USER_NOT_FOUND"
	ErrCodeProfileNotFound    = "PROFILE_NOT_FOUND"
)
