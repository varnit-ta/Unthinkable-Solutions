// Package handlers implements HTTP request handlers for the recipe API.
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/varnit-ta/smart-recipe-generator/backend/internal/auth"
	"github.com/varnit-ta/smart-recipe-generator/backend/internal/service"
)

// AuthHandler manages user authentication and registration endpoints.
type AuthHandler struct {
	Service   *service.Service
	JWTSecret string
	JWTExpiry int
}

// RegisterRequest contains user registration information.
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register handles POST /api/auth/register to create new user accounts.
//
// Request body: RegisterRequest with username, email, and password
//
// Security:
// - Password is hashed with bcrypt before storage
// - Returns JWT token on successful registration
//
// Returns: 200 OK with JWT token, or error status
func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "bad request"})
		return
	}

	user, err := a.Service.CreateUser(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "could not create user"})
		return
	}

	token, err := auth.GenerateJWT(a.JWTSecret, int(user.ID), a.JWTExpiry)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "could not generate token"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// LoginRequest contains user login credentials.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login handles POST /api/auth/login to authenticate existing users.
//
// Request body: LoginRequest with email and password
//
// Security:
// - Password is verified using bcrypt
// - Returns generic error message to prevent user enumeration
// - Issues JWT token on successful authentication
//
// Returns: 200 OK with JWT token, or 401 Unauthorized
func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "bad request"})
		return
	}

	user, err := a.Service.Authenticate(r.Context(), req.Email, req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "invalid credentials"})
		return
	}

	token, err := auth.GenerateJWT(a.JWTSecret, int(user.ID), a.JWTExpiry)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "could not generate token"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
}
