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

type Handler struct {
	Service       *service.Service
	VisionService vision.VisionService
	MaxImageBytes int64
}

func New(s *service.Service, vs vision.VisionService, maxImageMB int) *Handler {
	return &Handler{
		Service:       s,
		VisionService: vs,
		MaxImageBytes: int64(maxImageMB) * 1024 * 1024,
	}
}

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
		// Log the actual error
		println("SearchAndFilterRecipes error:", err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "database error"})
		return
	}
	// Convert to response format
	response := make([]RecipeDetailResponse, len(recipes))
	for i, r := range recipes {
		response[i] = toSearchRecipeResponse(r)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
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
	response := toRecipeDetailResponse(recipe)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
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
	// Convert to response format
	response := make([]FavoriteRecipeResponse, len(list))
	for i, fav := range list {
		response[i] = toFavoriteRecipeResponse(fav)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

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

// DetectIngredients accepts an image and uses vision AI to detect ingredients
func (h *Handler) DetectIngredients(w http.ResponseWriter, r *http.Request) {
	// Check if vision service is available
	if h.VisionService == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message":             "vision service not configured",
			"detectedIngredients": []string{},
		})
		return
	}

	// Parse multipart form with size limit
	if err := r.ParseMultipartForm(h.MaxImageBytes); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": "image too large or invalid form data",
		})
		return
	}

	// Get image file
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

	// Validate file type
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

	// Read image data
	imageData, err := io.ReadAll(file)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": "failed to read image",
		})
		return
	}

	// Check image size
	if len(imageData) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": "empty image file",
		})
		return
	}

	// Detect ingredients using vision service
	result, err := h.VisionService.DetectIngredients(r.Context(), imageData, filename)
	if err != nil {
		// Log error but return graceful response
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

	// Return successful response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"detectedIngredients": result.Ingredients,
		"confidence":          result.Confidence,
		"provider":            result.Provider,
		"caption":             result.RawResponse,
	})
}

// isValidImageType checks if the content type is a valid image format
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
