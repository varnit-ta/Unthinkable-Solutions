// Package middleware provides HTTP middleware functions for request processing.
package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logging is a middleware that logs HTTP requests with method, path, and duration.
// Logs are written in the format: "METHOD PATH DURATION"
// Example: "GET /recipes/123 15.2ms"
//
// This middleware should be applied globally to log all incoming requests.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
