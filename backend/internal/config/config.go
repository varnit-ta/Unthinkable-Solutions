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
	DatabaseURL    string
	Port           string
	JWTSecret      string
	JWTExpiryHours int
	DBMaxOpenConns int
	DBMaxIdleConns int
	DBConnMaxIdle  time.Duration
	DBConnMaxLife  time.Duration
	DBRetryMax     int
	DBRetryBackoff time.Duration
	AIServiceURL   string
	MaxImageSizeMB int
	AllowedOrigins string
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

	aiServiceURL := os.Getenv("AI_SERVICE_URL")
	if aiServiceURL == "" {
		aiServiceURL = "http://localhost:8000"
	}

	maxImageSize := parseIntEnv("MAX_IMAGE_SIZE_MB", 10)

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:5173,http://localhost:3000,http://localhost:4173,https://unthinkable-solutions-three.vercel.app/"
	}

	return Config{
		DatabaseURL:    db,
		Port:           port,
		JWTSecret:      secret,
		JWTExpiryHours: expiry,
		DBMaxOpenConns: maxOpen,
		DBMaxIdleConns: maxIdle,
		DBConnMaxIdle:  idle,
		DBConnMaxLife:  life,
		DBRetryMax:     retryMax,
		DBRetryBackoff: retryBackoff,
		AIServiceURL:   aiServiceURL,
		MaxImageSizeMB: maxImageSize,
		AllowedOrigins: allowedOrigins,
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
