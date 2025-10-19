// Package service implements the business logic layer for the recipe application.
// It provides recipe matching, search, filtering, user management, and recommendation algorithms.
package service

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strings"

	"github.com/varnit-ta/smart-recipe-generator/backend/internal/auth"
	"github.com/varnit-ta/smart-recipe-generator/backend/internal/db"
)

// Service provides business logic operations for the recipe application.
// It wraps the database queries and implements complex operations like scoring,
// filtering, and recommendations.
type Service struct {
	q *db.Queries
}

// NewService creates a new Service instance with the provided database connection.
//
// Parameters:
//   - conn: database connection implementing the DBTX interface
//
// Returns a Service ready to perform business operations.
func NewService(conn db.DBTX) *Service {
	return &Service{q: db.New(conn)}
}

// RecipeSummary represents a recipe with its match score.
// Used for ingredient-based recipe matching results.
type RecipeSummary struct {
	ID    int32  `json:"id"`
	Title string `json:"title"`
	Score int    `json:"score"`
}

// ListRecipes retrieves a paginated list of recipes.
//
// Parameters:
//   - ctx: request context
//   - limit: maximum number of recipes to return
//   - offset: number of recipes to skip (for pagination)
//
// Returns a slice of recipe rows or error.
func (s *Service) ListRecipes(ctx context.Context, limit, offset int) ([]db.ListRecipesRow, error) {
	params := db.ListRecipesParams{Limit: int32(limit), Offset: int32(offset)}
	return s.q.ListRecipes(ctx, params)
}

// GetRecipe retrieves complete recipe details by ID.
//
// Parameters:
//   - ctx: request context
//   - id: recipe identifier
//
// Returns full recipe data including ingredients, instructions, tags, and metadata.
func (s *Service) GetRecipe(ctx context.Context, id int) (db.GetRecipeByIDRow, error) {
	return s.q.GetRecipeByID(ctx, int32(id))
}

// MatchRecipes scores recipes based on ingredient overlap with detected items.
//
// Scoring algorithm:
// - +1 point for each detected ingredient matching a recipe tag
// - +1 point for each detected ingredient appearing in recipe title
// - Results sorted by descending score
//
// Parameters:
//   - ctx: request context
//   - detected: list of detected ingredient names
//   - limit: maximum recipes to check
//   - offset: pagination offset
//
// Returns scored recipe summaries sorted by relevance.
func (s *Service) MatchRecipes(ctx context.Context, detected []string, limit, offset int) ([]RecipeSummary, error) {
	list, err := s.ListRecipes(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	detectedSet := map[string]struct{}{}
	for _, d := range detected {
		detectedSet[strings.ToLower(strings.TrimSpace(d))] = struct{}{}
	}

	var results []RecipeSummary
	for _, r := range list {
		full, err := s.q.GetRecipeByID(ctx, r.ID)
		if err != nil {
			return nil, err
		}

		score := 0
		for _, t := range full.Tags {
			if _, ok := detectedSet[strings.ToLower(t)]; ok {
				score++
			}
		}

		titleLower := strings.ToLower(full.Title)
		for d := range detectedSet {
			if strings.Contains(titleLower, d) {
				score++
			}
		}

		results = append(results, RecipeSummary{ID: full.ID, Title: full.Title, Score: score})
	}

	sort.Slice(results, func(i, j int) bool { return results[i].Score > results[j].Score })
	return results, nil
}

// SearchAndFilterRecipes searches recipes and applies multiple optional filters.
//
// Filter behavior:
// - query: searches in recipe title and tags (empty = all recipes)
// - diet: matches recipe tags (e.g., "vegetarian", "vegan")
// - difficulty: exact match on difficulty level ("easy", "medium", "hard")
// - maxTimeMinutes: filters recipes by cooking time
// - cuisine: exact match on cuisine type
//
// Parameters:
//   - ctx: request context
//   - query: search query for title/tags
//   - diet: dietary restriction filter
//   - difficulty: difficulty level filter
//   - maxTimeMinutes: maximum cooking time in minutes (nil = no limit)
//   - cuisine: cuisine type filter
//   - limit: maximum results to return
//   - offset: pagination offset
//
// Returns filtered and paginated recipe list.
func (s *Service) SearchAndFilterRecipes(
	ctx context.Context,
	query string,
	diet string,
	difficulty string,
	maxTimeMinutes *int,
	cuisine string,
	limit int,
	offset int,
) ([]db.SearchRecipesRow, error) {
	fetchLimit := int32(math.Max(float64(limit+offset), 200))
	if fetchLimit > 2000 {
		fetchLimit = 2000
	}

	params := db.SearchRecipesParams{Column1: sql.NullString{String: query, Valid: true}, Limit: fetchLimit, Offset: 0}
	all, err := s.q.SearchRecipes(ctx, params)
	if err != nil {
		return nil, err
	}

	dietLower := strings.ToLower(strings.TrimSpace(diet))
	diffLower := strings.ToLower(strings.TrimSpace(difficulty))
	cuisineLower := strings.ToLower(strings.TrimSpace(cuisine))

	var filtered []db.SearchRecipesRow
	for _, r := range all {
		if diffLower != "" {
			if !r.Difficulty.Valid || strings.ToLower(r.Difficulty.String) != diffLower {
				continue
			}
		}
		if cuisineLower != "" {
			if !r.Cuisine.Valid || strings.ToLower(r.Cuisine.String) != cuisineLower {
				continue
			}
		}
		if maxTimeMinutes != nil {
			if !r.CookTimeMinutes.Valid || int(r.CookTimeMinutes.Int32) > *maxTimeMinutes {
				continue
			}
		}
		if dietLower != "" {
			matched := false
			for _, t := range r.Tags {
				if strings.ToLower(t) == dietLower {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}
		filtered = append(filtered, r)
	}

	if offset >= len(filtered) {
		return []db.SearchRecipesRow{}, nil
	}

	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[offset:end], nil
}

// MatchFilters defines optional filters for ingredient-based recipe matching.
type MatchFilters struct {
	Diet           string
	Difficulty     string
	MaxTimeMinutes *int
	Cuisine        string
	Limit          int
	Offset         int
}

// RecipeWithScore extends a recipe search result with a relevance score.
// Used for filtered matching operations.
type RecipeWithScore struct {
	db.SearchRecipesRow
	Score int `json:"score"`
}

// MatchWithFilters combines filtering and ingredient-based scoring.
//
// Process:
// 1. Apply all filters (diet, difficulty, time, cuisine)
// 2. Score remaining recipes by ingredient overlap
// 3. Sort by descending score
//
// Parameters:
//   - ctx: request context
//   - ingredients: list of ingredient names to match
//   - filters: optional filters to narrow results
//
// Returns scored and sorted recipes matching all criteria.
func (s *Service) MatchWithFilters(ctx context.Context, ingredients []string, filters MatchFilters) ([]RecipeWithScore, error) {
	candidates, err := s.SearchAndFilterRecipes(ctx, "", filters.Diet, filters.Difficulty, filters.MaxTimeMinutes, filters.Cuisine, filters.Limit, filters.Offset)
	if err != nil {
		return nil, err
	}
	detectedSet := map[string]struct{}{}
	for _, d := range ingredients {
		detectedSet[strings.ToLower(strings.TrimSpace(d))] = struct{}{}
	}

	var results []RecipeWithScore
	for _, r := range candidates {
		score := 0
		for _, t := range r.Tags {
			if _, ok := detectedSet[strings.ToLower(t)]; ok {
				score++
			}
		}

		titleLower := strings.ToLower(r.Title)
		for d := range detectedSet {
			if strings.Contains(titleLower, d) {
				score++
			}
		}

		results = append(results, RecipeWithScore{SearchRecipesRow: r, Score: score})
	}

	sort.Slice(results, func(i, j int) bool { return results[i].Score > results[j].Score })
	return results, nil
}

// CreateUser registers a new user with hashed password.
//
// Security:
// - Password is hashed using bcrypt before storage
// - Original password is never stored
//
// Parameters:
//   - ctx: request context
//   - username: user's chosen username
//   - email: user's email address
//   - password: plain text password (will be hashed)
//
// Returns created user data or error if registration fails.
func (s *Service) CreateUser(ctx context.Context, username, email, password string) (db.CreateUserRow, error) {
	hash, err := auth.HashPassword(password)
	if err != nil {
		return db.CreateUserRow{}, err
	}

	params := db.CreateUserParams{
		Username:     sql.NullString{String: username, Valid: true},
		Email:        sql.NullString{String: email, Valid: true},
		PasswordHash: sql.NullString{String: hash, Valid: true},
	}
	return s.q.CreateUser(ctx, params)
}

// Authenticate verifies user credentials for login.
//
// Security:
// - Uses bcrypt for password verification
// - Returns generic error on failure (no user enumeration)
//
// Parameters:
//   - ctx: request context
//   - email: user's email address
//   - password: plain text password to verify
//
// Returns user data on success, error on authentication failure.
func (s *Service) Authenticate(ctx context.Context, email, password string) (db.GetUserByEmailRow, error) {
	row, err := s.q.GetUserByEmail(ctx, sql.NullString{String: email, Valid: true})
	if err != nil {
		return db.GetUserByEmailRow{}, fmt.Errorf("auth failed")
	}

	if err := auth.VerifyPassword(row.PasswordHash.String, password); err != nil {
		return db.GetUserByEmailRow{}, fmt.Errorf("auth failed")
	}

	return row, nil
}

// AddRating records a user's rating for a recipe.
//
// Parameters:
//   - ctx: request context
//   - userID: ID of the user submitting the rating
//   - recipeID: ID of the recipe being rated
//   - rating: numeric rating value
//
// Returns the created rating record or error.
func (s *Service) AddRating(ctx context.Context, userID sql.NullInt32, recipeID, rating int) (db.Rating, error) {
	params := db.InsertRatingParams{UserID: userID, RecipeID: sql.NullInt32{Int32: int32(recipeID), Valid: true}, Rating: sql.NullInt32{Int32: int32(rating), Valid: true}}
	return s.q.InsertRating(ctx, params)
}

// AddFavorite adds a recipe to a user's favorites list.
//
// Parameters:
//   - ctx: request context
//   - userID: ID of the user
//   - recipeID: ID of the recipe to favorite
//
// Returns the created favorite record or error.
func (s *Service) AddFavorite(ctx context.Context, userID int, recipeID int) (db.Favorite, error) {
	params := db.AddFavoriteParams{UserID: sql.NullInt32{Int32: int32(userID), Valid: true}, RecipeID: sql.NullInt32{Int32: int32(recipeID), Valid: true}}
	return s.q.AddFavorite(ctx, params)
}

// RemoveFavorite removes a recipe from a user's favorites list.
//
// Parameters:
//   - ctx: request context
//   - userID: ID of the user
//   - recipeID: ID of the recipe to unfavorite
//
// Returns error if operation fails.
func (s *Service) RemoveFavorite(ctx context.Context, userID int, recipeID int) error {
	params := db.RemoveFavoriteParams{UserID: sql.NullInt32{Int32: int32(userID), Valid: true}, RecipeID: sql.NullInt32{Int32: int32(recipeID), Valid: true}}
	return s.q.RemoveFavorite(ctx, params)
}

// ListFavorites retrieves all recipes favorited by a user.
//
// Parameters:
//   - ctx: request context
//   - userID: ID of the user
//
// Returns list of favorited recipes or error.
func (s *Service) ListFavorites(ctx context.Context, userID int) ([]db.ListFavoritesByUserRow, error) {
	return s.q.ListFavoritesByUser(ctx, sql.NullInt32{Int32: int32(userID), Valid: true})
}

// IsFavorite checks if a recipe is in a user's favorites.
//
// Parameters:
//   - ctx: request context
//   - userID: ID of the user
//   - recipeID: ID of the recipe to check
//
// Returns true if recipe is favorited, false otherwise.
func (s *Service) IsFavorite(ctx context.Context, userID int, recipeID int) (bool, error) {
	return s.q.IsFavorite(ctx, db.IsFavoriteParams{
		UserID:   sql.NullInt32{Int32: int32(userID), Valid: true},
		RecipeID: sql.NullInt32{Int32: int32(recipeID), Valid: true},
	})
}

// GetSuggestions generates personalized recipe recommendations for a user.
//
// Recommendation algorithm (content-based filtering):
// 1. Analyze user's favorite recipes
// 2. Extract and count tags from favorites
// 3. Score candidate recipes by tag overlap with favorites
// 4. Return top-scored recipes
//
// The algorithm favors recipes with tags that frequently appear in
// the user's favorites, creating personalized recommendations based on
// demonstrated preferences.
//
// Parameters:
//   - ctx: request context
//   - userID: ID of the user to generate suggestions for
//   - limit: maximum number of suggestions to return
//
// Returns scored recipe suggestions or error.
func (s *Service) GetSuggestions(ctx context.Context, userID int, limit int) ([]RecipeWithScore, error) {
	favs, err := s.ListFavorites(ctx, userID)
	if err != nil {
		return nil, err
	}

	favoriteTagCounts := map[string]int{}
	for _, f := range favs {
		full, err := s.q.GetRecipeByID(ctx, f.RecipeID.Int32)
		if err != nil {
			continue
		}
		for _, t := range full.Tags {
			favoriteTagCounts[strings.ToLower(t)]++
		}
	}

	candidates, err := s.SearchAndFilterRecipes(ctx, "", "", "", nil, "", int(math.Max(float64(limit*5), 100)), 0)
	if err != nil {
		return nil, err
	}

	var scored []RecipeWithScore
	for _, c := range candidates {
		score := 0
		for _, t := range c.Tags {
			score += favoriteTagCounts[strings.ToLower(t)]
		}
		if score > 0 {
			scored = append(scored, RecipeWithScore{SearchRecipesRow: c, Score: score})
		}
	}

	sort.Slice(scored, func(i, j int) bool { return scored[i].Score > scored[j].Score })
	if len(scored) > limit {
		scored = scored[:limit]
	}

	return scored, nil
}

// ErrBadRequest is a sentinel error for invalid requests.
var (
	ErrBadRequest = fmt.Errorf("%d", http.StatusBadRequest)
)
