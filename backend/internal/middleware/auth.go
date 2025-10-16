// Package middleware provides HTTP middleware functions for request processing.
// This includes authentication, logging, and other cross-cutting concerns.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/varnit-ta/smart-recipe-generator/backend/internal/auth"
)

// ctxKey is a custom type for context keys to avoid collisions.
type ctxKey string

// UserIDKey is the context key used to store the authenticated user's ID.
// Handlers can retrieve the user ID from the request context using this key.
const UserIDKey ctxKey = "userId"

// JWTAuth returns a middleware function that validates JWT tokens.
// It extracts the token from the Authorization header (format: "Bearer <token>"),
// validates it, and stores the user ID in the request context.
//
// Protected routes should use this middleware to ensure authentication.
// If authentication fails, returns 401 Unauthorized.
//
// Parameters:
//   - secret: The secret key used to verify JWT token signatures
//
// Returns a middleware function that can be chained with Chi router.
func JWTAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]
			claims, err := auth.ParseJWT(secret, tokenString)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
