package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/varnit-ta/smart-recipe-generator/backend/internal/auth"
	"github.com/varnit-ta/smart-recipe-generator/backend/internal/service"
)

type AuthHandler struct {
	Service   *service.Service
	JWTSecret string
	JWTExpiry int
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", 400)
		return
	}
	// create user via service
	user, err := a.Service.CreateUser(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		http.Error(w, "could not create user", 500)
		return
	}
	token, err := auth.GenerateJWT(a.JWTSecret, int(user.ID), a.JWTExpiry)
	if err != nil {
		http.Error(w, "could not generate token", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", 400)
		return
	}
	user, err := a.Service.Authenticate(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	token, err := auth.GenerateJWT(a.JWTSecret, int(user.ID), a.JWTExpiry)
	if err != nil {
		http.Error(w, "could not generate token", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
}
