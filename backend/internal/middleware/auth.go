package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/varnit-ta/smart-recipe-generator/backend/internal/auth"
)

type ctxKey string

const UserIDKey ctxKey = "userId"

func JWTAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hdr := r.Header.Get("Authorization")
			if hdr == "" {
				http.Error(w, "unauthorized", 401)
				return
			}
			parts := strings.SplitN(hdr, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "unauthorized", 401)
				return
			}
			tokenStr := parts[1]
			claims, err := auth.ParseJWT(secret, tokenStr)
			if err != nil {
				http.Error(w, "unauthorized", 401)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
