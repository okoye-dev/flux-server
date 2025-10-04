package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for our application
type Config struct {
	Server   ServerConfig
	Supabase SupabaseConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port        string
	Environment string
}

// SupabaseConfig holds Supabase-related configuration
type SupabaseConfig struct {
	URL             string
	AnonKey         string
	ServiceRoleKey  string
	JWTSecret       string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:        getEnv("PORT", "8080"),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
		Supabase: SupabaseConfig{
			URL:            getEnv("SUPABASE_URL", ""),
			AnonKey:        getEnv("SUPABASE_ANON_KEY", ""),
			ServiceRoleKey: getEnv("SUPABASE_SERVICE_ROLE_KEY", ""),
			JWTSecret:      getEnv("JWT_SECRET", ""),
		},
	}
}

// Validate checks if required configuration is present
func (c *Config) Validate() error {
	if c.Supabase.URL == "" {
		return &ConfigError{Field: "SUPABASE_URL", Message: "Supabase URL is required"}
	}
	if c.Supabase.AnonKey == "" {
		return &ConfigError{Field: "SUPABASE_ANON_KEY", Message: "Supabase anon key is required"}
	}
	if c.Supabase.JWTSecret == "" {
		return &ConfigError{Field: "JWT_SECRET", Message: "JWT secret is required for secure token validation"}
	}
	return nil
}

// ConfigError represents a configuration error
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return e.Message
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets an environment variable as integer with a fallback value
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}
