// Package config provides configuration management for the application.
// It loads settings from environment variables with sensible defaults.
package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration settings loaded from environment variables.
// Each field has a corresponding environment variable and default value.
type Config struct {
	DatabaseURL       string        // PostgreSQL connection string (DATABASE_URL)
	Port              string        // HTTP server port (PORT, default: 8081)
	JWTSecret         string        // Secret key for JWT signing (JWT_SECRET)
	JWTExpiryHours    int           // JWT token expiration in hours (default: 48)
	DBMaxOpenConns    int           // Maximum open database connections (DB_MAX_OPEN_CONNS, default: 20)
	DBMaxIdleConns    int           // Maximum idle database connections (DB_MAX_IDLE_CONNS, default: 10)
	DBConnMaxIdle     time.Duration // Maximum connection idle time (DB_CONN_MAX_IDLE, default: 15m)
	DBConnMaxLife     time.Duration // Maximum connection lifetime (DB_CONN_MAX_LIFE, default: 1h)
	DBRetryMax        int           // Maximum database connection retry attempts (DB_RETRY_MAX, default: 8)
	DBRetryBackoff    time.Duration // Initial retry backoff duration (DB_RETRY_BACKOFF, default: 500ms)
	HuggingFaceAPIKey string        // Hugging Face API key for vision AI (HUGGINGFACE_API_KEY)
	MaxImageSizeMB    int           // Maximum image upload size in MB (MAX_IMAGE_SIZE_MB, default: 10)
	AllowedOrigins    string        // Comma-separated list of allowed CORS origins (ALLOWED_ORIGINS)
}

// Load reads configuration from environment variables and returns a Config struct
// with all settings initialized. Provides sensible defaults for all optional settings.
func Load() Config {
	db := os.Getenv("DATABASE_URL")
	if db == "" {
		db = "postgres://unthinkable:unthinkable@localhost:5432/unthinkable_recipes?sslmode=disable"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "change-me-to-a-secure-secret"
	}
	expiry := 48

	maxOpen := parseIntEnv("DB_MAX_OPEN_CONNS", 20)
	maxIdle := parseIntEnv("DB_MAX_IDLE_CONNS", 10)
	idle := parseDurationEnv("DB_CONN_MAX_IDLE", 15*time.Minute)
	life := parseDurationEnv("DB_CONN_MAX_LIFE", time.Hour)
	retryMax := parseIntEnv("DB_RETRY_MAX", 8)
	retryBackoff := parseDurationEnv("DB_RETRY_BACKOFF", 500*time.Millisecond)

	hfAPIKey := os.Getenv("HUGGINGFACE_API_KEY")
	maxImageSize := parseIntEnv("MAX_IMAGE_SIZE_MB", 10)

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:5173,http://localhost:3000,http://localhost:4173,https://unthinkable-solutions-three.vercel.app/"
	}

	return Config{
		DatabaseURL:       db,
		Port:              port,
		JWTSecret:         secret,
		JWTExpiryHours:    expiry,
		DBMaxOpenConns:    maxOpen,
		DBMaxIdleConns:    maxIdle,
		DBConnMaxIdle:     idle,
		DBConnMaxLife:     life,
		DBRetryMax:        retryMax,
		DBRetryBackoff:    retryBackoff,
		HuggingFaceAPIKey: hfAPIKey,
		MaxImageSizeMB:    maxImageSize,
		AllowedOrigins:    allowedOrigins,
	}
}

// parseIntEnv reads an integer environment variable with a default fallback.
// Returns the default value if the variable is not set or cannot be parsed.
func parseIntEnv(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	if n, err := strconv.Atoi(v); err == nil {
		return n
	}
	return def
}

// parseDurationEnv reads a duration environment variable with a default fallback.
// Supports Go duration format (e.g., "30s", "5m", "2h").
// Returns the default value if the variable is not set or cannot be parsed.
func parseDurationEnv(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	if d, err := time.ParseDuration(v); err == nil {
		return d
	}
	return def
}
