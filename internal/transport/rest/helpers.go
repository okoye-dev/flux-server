package rest

import (
	"encoding/json"
	"net/http"
	"time"
)

// HTTP Response Helpers

// WriteSuccessResponse writes a successful API response
func WriteSuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	response := APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// WriteErrorResponse writes an error API response
func WriteErrorResponse(w http.ResponseWriter, statusCode int, errorCode, message, details string) {
	response := APIResponse{
		Success: false,
		Message: message,
		Error: &APIError{
			Code:    errorCode,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// WriteValidationErrorResponse writes a validation error response
func WriteValidationErrorResponse(w http.ResponseWriter, errors []ValidationError) {
	response := APIResponse{
		Success: false,
		Message: MsgInvalidRequest,
		Error: &APIError{
			Code:    ErrCodeValidation,
			Message: MsgInvalidRequest,
			Details: "Validation failed",
		},
		Data:      ValidationErrorResponse{Errors: errors},
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(response)
}

// WritePaginatedResponse writes a paginated response
func WritePaginatedResponse(w http.ResponseWriter, data interface{}, pagination Pagination) {
	response := PaginatedResponse{
		Data:       data,
		Pagination: pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Specific Response Helpers

// WriteAuthResponse writes an authentication response
func WriteAuthResponse(w http.ResponseWriter, statusCode int, authResp AuthResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(authResp)
}

// WriteProfileResponse writes a profile response
func WriteProfileResponse(w http.ResponseWriter, profileResp ProfileResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profileResp)
}

// WriteProtectedDataResponse writes a protected data response
func WriteProtectedDataResponse(w http.ResponseWriter, protectedResp ProtectedDataResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(protectedResp)
}

// WriteHealthResponse writes a health check response
func WriteHealthResponse(w http.ResponseWriter, healthResp HealthResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(healthResp)
}

// WriteRootResponse writes a root endpoint response
func WriteRootResponse(w http.ResponseWriter, rootResp RootResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rootResp)
}

// Common Error Helpers

// WriteUnauthorizedError writes an unauthorized error response
func WriteUnauthorizedError(w http.ResponseWriter, message string) {
	if message == "" {
		message = MsgUnauthorized
	}
	WriteErrorResponse(w, http.StatusUnauthorized, ErrCodeUnauthorized, message, "")
}

// WriteForbiddenError writes a forbidden error response
func WriteForbiddenError(w http.ResponseWriter, message string) {
	if message == "" {
		message = MsgForbidden
	}
	WriteErrorResponse(w, http.StatusForbidden, ErrCodeForbidden, message, "")
}

// WriteNotFoundError writes a not found error response
func WriteNotFoundError(w http.ResponseWriter, message string) {
	if message == "" {
		message = MsgNotFound
	}
	WriteErrorResponse(w, http.StatusNotFound, ErrCodeNotFound, message, "")
}

// WriteInternalServerError writes an internal server error response
func WriteInternalServerError(w http.ResponseWriter, message string, details string) {
	if message == "" {
		message = MsgInternalServerError
	}
	WriteErrorResponse(w, http.StatusInternalServerError, ErrCodeInternalError, message, details)
}

// WriteBadRequestError writes a bad request error response
func WriteBadRequestError(w http.ResponseWriter, message string, details string) {
	if message == "" {
		message = MsgInvalidRequest
	}
	WriteErrorResponse(w, http.StatusBadRequest, ErrCodeValidation, message, details)
}

// WriteMethodNotAllowedError writes a method not allowed error response
func WriteMethodNotAllowedError(w http.ResponseWriter) {
	WriteErrorResponse(w, http.StatusMethodNotAllowed, ErrCodeValidation, "Method not allowed", "")
}

// WriteSupabaseError writes a Supabase error response
func WriteSupabaseError(w http.ResponseWriter, message string, details string) {
	WriteErrorResponse(w, http.StatusBadRequest, ErrCodeSupabaseError, message, details)
}

// WriteAuthError writes an authentication error response
func WriteAuthError(w http.ResponseWriter, message string, details string) {
	WriteErrorResponse(w, http.StatusUnauthorized, ErrCodeAuthError, message, details)
}

// WriteInvalidTokenError writes an invalid token error response
func WriteInvalidTokenError(w http.ResponseWriter) {
	WriteErrorResponse(w, http.StatusUnauthorized, ErrCodeInvalidToken, MsgInvalidOrExpiredToken, "")
}

// WriteMissingConfigError writes a missing configuration error response
func WriteMissingConfigError(w http.ResponseWriter, config string) {
	WriteErrorResponse(w, http.StatusInternalServerError, ErrCodeMissingConfig, MsgSupabaseConfigMissing, config)
}
