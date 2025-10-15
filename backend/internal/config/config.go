package config

import (
	"os"
	"strconv"
	"time"
)

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
}

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
	// pooling defaults
	maxOpen := parseIntEnv("DB_MAX_OPEN_CONNS", 20)
	maxIdle := parseIntEnv("DB_MAX_IDLE_CONNS", 10)
	idle := parseDurationEnv("DB_CONN_MAX_IDLE", time.Minute*15)
	life := parseDurationEnv("DB_CONN_MAX_LIFE", time.Hour)
	retryMax := parseIntEnv("DB_RETRY_MAX", 8)
	retryBackoff := parseDurationEnv("DB_RETRY_BACKOFF", 500*time.Millisecond) // default 500ms

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
	}
}

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
