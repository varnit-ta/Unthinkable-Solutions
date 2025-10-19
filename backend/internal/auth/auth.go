// Package auth provides authentication and authorization utilities.
// It handles JWT token generation/parsing and password hashing/verification
// using industry-standard cryptographic libraries.
package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Claims represents the JWT token payload containing user identification
// and standard JWT claims (expiration, issued at, etc.).
type Claims struct {
	UserID int `json:"userId"`
	jwt.RegisteredClaims
}

// HashPassword generates a bcrypt hash of the provided password.
// Uses bcrypt.DefaultCost (currently 10) for the hashing cost factor.
// This is computationally expensive by design to resist brute-force attacks.
//
// Returns the hashed password string or an error if hashing fails.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// VerifyPassword compares a bcrypt hash with a plain-text password.
// Returns nil if the password matches the hash, otherwise returns an error.
//
// This is a constant-time comparison that prevents timing attacks.
func VerifyPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// GenerateJWT creates a signed JWT token for the given user.
// The token includes the user ID in the claims and is signed using HMAC-SHA256.
//
// Parameters:
//   - secret: Secret key for signing the token
//   - userID: User identifier to embed in the token
//   - expiryHours: Number of hours until token expiration
//
// Returns the signed token string or an error if signing fails.
func GenerateJWT(secret string, userID int, expiryHours int) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(expiryHours) * time.Hour)

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseJWT validates and parses a JWT token string.
// Verifies the signature, expiration, and extracts the claims.
//
// Parameters:
//   - secret: Secret key used to verify the token signature
//   - tokenStr: JWT token string to parse
//
// Returns the parsed claims or an error if validation fails.
// Common errors include expired tokens, invalid signatures, or malformed tokens.
func ParseJWT(secret, tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// RandomSecret generates a cryptographically secure random secret key.
// Generates 32 random bytes and encodes them as base64.
// Suitable for use as JWT signing secrets or other security tokens.
//
// Returns a base64-encoded random string or an error if random generation fails.
func RandomSecret() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
