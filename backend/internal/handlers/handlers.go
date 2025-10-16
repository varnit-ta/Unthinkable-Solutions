package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/varnit-ta/smart-recipe-generator/backend/internal/middleware"
	"github.com/varnit-ta/smart-recipe-generator/backend/internal/service"
)

type Handler struct {
	Service *service.Service
}

func New(s *service.Service) *Handler { return &Handler{Service: s} }

func (h *Handler) ListRecipes(w http.ResponseWriter, r *http.Request) {
	// Parse query params
	q := r.URL.Query().Get("q")
	diet := r.URL.Query().Get("diet")
	difficulty := r.URL.Query().Get("difficulty")
	cuisine := r.URL.Query().Get("cuisine")
	maxTimeStr := r.URL.Query().Get("maxTime")
	limit := 50
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 200 {
			limit = n
		}
	}
	offset := 0
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}
	var maxTimePtr *int
	if maxTimeStr != "" {
		if n, err := strconv.Atoi(maxTimeStr); err == nil && n > 0 {
			maxTimePtr = &n
		}
	}
	recipes, err := h.Service.SearchAndFilterRecipes(r.Context(), q, diet, difficulty, maxTimePtr, cuisine, limit, offset)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "database error"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(recipes)
}

func (h *Handler) GetRecipe(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	recipe, err := h.Service.GetRecipe(r.Context(), id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "recipe not found"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(recipe)
}

type MatchRequest struct {
	DetectedIngredients []string `json:"detectedIngredients"`
}

func (h *Handler) Match(w http.ResponseWriter, r *http.Request) {
	var req MatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "bad request"})
		return
	}
	// optional filters via query
	diet := r.URL.Query().Get("diet")
	difficulty := r.URL.Query().Get("difficulty")
	cuisine := r.URL.Query().Get("cuisine")
	maxTimeStr := r.URL.Query().Get("maxTime")
	limit := 50
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 200 {
			limit = n
		}
	}
	offset := 0
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}
	var maxTimePtr *int
	if maxTimeStr != "" {
		if n, err := strconv.Atoi(maxTimeStr); err == nil && n > 0 {
			maxTimePtr = &n
		}
	}
	recipes, err := h.Service.MatchWithFilters(r.Context(), req.DetectedIngredients, service.MatchFilters{
		Diet: diet, Difficulty: difficulty, MaxTimeMinutes: maxTimePtr, Cuisine: cuisine, Limit: limit, Offset: offset,
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "server error"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(recipes)
}

type RatingRequest struct {
	RecipeID int `json:"recipeId"`
	Rating   int `json:"rating"`
}

func (h *Handler) PostRating(w http.ResponseWriter, r *http.Request) {
	var req RatingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "bad request"})
		return
	}
	// extract user id from context populated by JWTAuth
	var uid32 sql.NullInt32
	if v := r.Context().Value(middleware.UserIDKey); v != nil {
		if id, ok := v.(int); ok {
			uid32 = sql.NullInt32{Int32: int32(id), Valid: true}
		}
	}
	rt, err := h.Service.AddRating(r.Context(), uid32, req.RecipeID, req.Rating)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "server error"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(rt)
}

func (h *Handler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	// URL param recipeId
	idStr := chi.URLParam(r, "id")
	recipeID, err := strconv.Atoi(idStr)
	if err != nil || recipeID <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "bad request"})
		return
	}
	v := r.Context().Value(middleware.UserIDKey)
	id, ok := v.(int)
	if !ok || id <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "unauthorized"})
		return
	}
	fav, err := h.Service.AddFavorite(r.Context(), id, recipeID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "server error"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(fav)
}

func (h *Handler) RemoveFavorite(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	recipeID, err := strconv.Atoi(idStr)
	if err != nil || recipeID <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "bad request"})
		return
	}
	v := r.Context().Value(middleware.UserIDKey)
	id, ok := v.(int)
	if !ok || id <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "unauthorized"})
		return
	}
	if err := h.Service.RemoveFavorite(r.Context(), id, recipeID); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "server error"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListFavorites(w http.ResponseWriter, r *http.Request) {
	v := r.Context().Value(middleware.UserIDKey)
	id, ok := v.(int)
	if !ok || id <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "unauthorized"})
		return
	}
	list, err := h.Service.ListFavorites(r.Context(), id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "server error"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(list)
}

// DetectIngredients is a stub that accepts an image and returns mocked detections.
func (h *Handler) DetectIngredients(w http.ResponseWriter, r *http.Request) {
	// Accept multipart form with key "image"
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB
		http.Error(w, "bad request", 400)
		return
	}
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "bad request", 400)
		return
	}
	_ = file.Close()
	name := ""
	if header != nil {
		name = strings.ToLower(header.Filename)
	}
	detected := []string{}
	// naive string contains to mock detections
	keywords := []string{"tomato", "chicken", "onion", "garlic", "egg", "cheese", "banana", "pepper", "rice", "pasta"}
	for _, k := range keywords {
		if strings.Contains(name, k) {
			detected = append(detected, k)
		}
	}
	if len(detected) == 0 {
		// default stub outputs
		detected = []string{"onion", "salt"}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"detectedIngredients": detected})
}

func (h *Handler) GetSuggestions(w http.ResponseWriter, r *http.Request) {
	v := r.Context().Value(middleware.UserIDKey)
	id, ok := v.(int)
	if !ok || id <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "unauthorized"})
		return
	}
	limit := 10
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}
	list, err := h.Service.GetSuggestions(r.Context(), id, limit)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "server error"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(list)
}
