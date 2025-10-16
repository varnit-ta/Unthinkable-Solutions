// Package handlers implements HTTP request handlers for the recipe API.
// Provides endpoints for recipes, matching, favorites, ratings, and AI-powered ingredient detection.
package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/varnit-ta/smart-recipe-generator/backend/internal/middleware"
	"github.com/varnit-ta/smart-recipe-generator/backend/internal/service"
	"github.com/varnit-ta/smart-recipe-generator/backend/internal/vision"
)

// Handler manages HTTP requests for recipe-related operations.
type Handler struct {
	Service       *service.Service     // Business logic service
	VisionService vision.VisionService // AI ingredient detection service
	MaxImageBytes int64                // Maximum upload size in bytes
}

// New creates a Handler with configured services and limits.
//
// Parameters:
//   - s: business logic service
//   - vs: vision AI service (can be nil to disable image detection)
//   - maxImageMB: maximum image upload size in megabytes
//
// Returns a configured Handler ready to serve requests.
func New(s *service.Service, vs vision.VisionService, maxImageMB int) *Handler {
	return &Handler{
		Service:       s,
		VisionService: vs,
		MaxImageBytes: int64(maxImageMB) * 1024 * 1024,
	}
}

// ListRecipes handles GET /api/recipes with search and filtering.
//
// Query parameters:
//   - q: search query for title/tags
//   - diet: dietary restriction (e.g., "vegetarian")
//   - difficulty: "easy", "medium", or "hard"
//   - cuisine: cuisine type filter
//   - maxTime: maximum cooking time in minutes
//   - limit: results per page (default 50, max 200)
//   - offset: pagination offset
//
// Returns: 200 OK with recipe array or error
func (h *Handler) ListRecipes(w http.ResponseWriter, r *http.Request) {
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
		println("SearchAndFilterRecipes error:", err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "database error"})
		return
	}

	response := make([]RecipeDetailResponse, len(recipes))
	for i, r := range recipes {
		response[i] = toSearchRecipeResponse(r)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// GetRecipe handles GET /api/recipes/:id to retrieve full recipe details.
//
// Path parameters:
//   - id: recipe identifier
//
// Returns: 200 OK with recipe details, or 404 if not found
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

	response := toRecipeDetailResponse(recipe)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// MatchRequest contains ingredients detected from image analysis.
type MatchRequest struct {
	DetectedIngredients []string `json:"detectedIngredients"` // List of ingredient names
}

// Match handles POST /api/match to find recipes matching ingredients.
//
// Request body: MatchRequest with detectedIngredients array
// Query parameters: same as ListRecipes (diet, difficulty, etc.)
//
// Returns: 200 OK with scored recipes sorted by match relevance
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
	// Convert to response format
	type RecipeWithScoreResponse struct {
		RecipeDetailResponse
		Score int `json:"score"`
	}
	response := make([]RecipeWithScoreResponse, len(recipes))
	for i, r := range recipes {
		response[i] = RecipeWithScoreResponse{
			RecipeDetailResponse: toSearchRecipeResponse(r.SearchRecipesRow),
			Score:                r.Score,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// RatingRequest contains a user's recipe rating submission.
type RatingRequest struct {
	RecipeID int `json:"recipeId"` // Recipe being rated
	Rating   int `json:"rating"`   // Numeric rating value
}

// PostRating handles POST /api/ratings (requires authentication).
//
// Request body: RatingRequest with recipeId and rating
//
// Returns: 200 OK with created rating record
func (h *Handler) PostRating(w http.ResponseWriter, r *http.Request) {
	var req RatingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "bad request"})
		return
	}

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

// AddFavorite handles POST /api/favorites/:id (requires authentication).
//
// Path parameters:
//   - id: recipe identifier to favorite
//
// Returns: 201 Created with favorite record
func (h *Handler) AddFavorite(w http.ResponseWriter, r *http.Request) {
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

// RemoveFavorite handles DELETE /api/favorites/:id (requires authentication).
//
// Path parameters:
//   - id: recipe identifier to unfavorite
//
// Returns: 204 No Content on success
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

// ListFavorites handles GET /api/favorites (requires authentication).
//
// Returns: 200 OK with array of user's favorited recipes
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

	response := make([]FavoriteRecipeResponse, len(list))
	for i, fav := range list {
		response[i] = toFavoriteRecipeResponse(fav)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// IsFavorite handles GET /api/favorites/:id/status (requires authentication).
//
// Path parameters:
//   - id: recipe identifier to check
//
// Returns: 200 OK with {"isFavorite": true/false}
func (h *Handler) IsFavorite(w http.ResponseWriter, r *http.Request) {
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

	isFav, err := h.Service.IsFavorite(r.Context(), id, recipeID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "server error"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]bool{"isFavorite": isFav})
}

// DetectIngredients handles POST /api/detect to extract ingredients from images.
//
// Request: multipart/form-data with "image" file field
// Supported formats: JPEG, PNG, GIF, WebP
// Max size: configured via MaxImageBytes
//
// Returns: 200 OK with detected ingredients and confidence score
func (h *Handler) DetectIngredients(w http.ResponseWriter, r *http.Request) {
	if h.VisionService == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message":             "vision service not configured",
			"detectedIngredients": []string{},
		})
		return
	}

	if err := r.ParseMultipartForm(h.MaxImageBytes); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": "image too large or invalid form data",
		})
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": "no image file provided",
		})
		return
	}
	defer file.Close()

	filename := ""
	if header != nil {
		filename = header.Filename
		contentType := header.Header.Get("Content-Type")
		if !isValidImageType(contentType) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"message": "invalid image format. Supported: JPEG, PNG, GIF, WebP",
			})
			return
		}
	}

	imageData, err := io.ReadAll(file)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": "failed to read image",
		})
		return
	}

	if len(imageData) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": "empty image file",
		})
		return
	}

	result, err := h.VisionService.DetectIngredients(r.Context(), imageData, filename)
	if err != nil {
		fmt.Printf("Vision API error: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"detectedIngredients": []string{},
			"message":             "Could not detect ingredients. Please try again or add them manually.",
			"error":               err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"detectedIngredients": result.Ingredients,
		"confidence":          result.Confidence,
		"provider":            result.Provider,
		"caption":             result.RawResponse,
	})
}

// isValidImageType validates that the uploaded file is a supported image format.
//
// Supported types: JPEG, PNG, GIF, WebP
func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
	}
	for _, t := range validTypes {
		if strings.Contains(strings.ToLower(contentType), t) {
			return true
		}
	}
	return false
}

// GetSuggestions handles GET /api/suggestions (requires authentication).
//
// Generates personalized recipe recommendations based on user's favorites.
//
// Query parameters:
//   - limit: maximum suggestions to return (default 10, max 100)
//
// Returns: 200 OK with scored recipe suggestions
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
	// Convert to response format
	type RecipeWithScoreResponse struct {
		RecipeDetailResponse
		Score int `json:"score"`
	}
	response := make([]RecipeWithScoreResponse, len(list))
	for i, r := range list {
		response[i] = RecipeWithScoreResponse{
			RecipeDetailResponse: toSearchRecipeResponse(r.SearchRecipesRow),
			Score:                r.Score,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
