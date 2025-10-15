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

type Service struct {
	q *db.Queries
}

func NewService(conn db.DBTX) *Service {
	return &Service{q: db.New(conn)}
}

type RecipeSummary struct {
	ID    int32  `json:"id"`
	Title string `json:"title"`
	Score int    `json:"score"`
}

// ListRecipes uses sqlc-generated ListRecipes
func (s *Service) ListRecipes(ctx context.Context, limit, offset int) ([]db.ListRecipesRow, error) {
	params := db.ListRecipesParams{Limit: int32(limit), Offset: int32(offset)}
	return s.q.ListRecipes(ctx, params)
}

// GetRecipe returns full recipe row
func (s *Service) GetRecipe(ctx context.Context, id int) (db.GetRecipeByIDRow, error) {
	return s.q.GetRecipeByID(ctx, int32(id))
}

// MatchRecipes: score by tag/title overlap
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
		// title contains matches
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

// SearchAndFilterRecipes performs a search by title/tags and applies optional filters.
// diet matches against tags; difficulty exact match; cuisine exact match; maxTime in minutes.
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
	// Use SearchRecipes with empty query to fetch all when query is blank.
	// Fetch more than needed to accommodate post-filtering; cap reasonably.
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
	// Apply pagination slice
	if offset >= len(filtered) {
		return []db.SearchRecipesRow{}, nil
	}
	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[offset:end], nil
}

type MatchFilters struct {
	Diet           string
	Difficulty     string
	MaxTimeMinutes *int
	Cuisine        string
	Limit          int
	Offset         int
}

// MatchWithFilters filters candidates first, then scores them by title/tag overlap against provided ingredients.
func (s *Service) MatchWithFilters(ctx context.Context, ingredients []string, filters MatchFilters) ([]RecipeSummary, error) {
	candidates, err := s.SearchAndFilterRecipes(ctx, "", filters.Diet, filters.Difficulty, filters.MaxTimeMinutes, filters.Cuisine, filters.Limit, filters.Offset)
	if err != nil {
		return nil, err
	}
	detectedSet := map[string]struct{}{}
	for _, d := range ingredients {
		detectedSet[strings.ToLower(strings.TrimSpace(d))] = struct{}{}
	}
	var results []RecipeSummary
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
		results = append(results, RecipeSummary{ID: r.ID, Title: r.Title, Score: score})
	}
	sort.Slice(results, func(i, j int) bool { return results[i].Score > results[j].Score })
	return results, nil
}

// CreateUser uses sqlc CreateUser
func (s *Service) CreateUser(ctx context.Context, username, email, password string) (db.CreateUserRow, error) {
	hash, err := auth.HashPassword(password)
	if err != nil {
		return db.CreateUserRow{}, err
	}
	params := db.CreateUserParams{Username: sql.NullString{String: username, Valid: true}, Email: sql.NullString{String: email, Valid: true}, PasswordHash: sql.NullString{String: hash, Valid: true}}
	return s.q.CreateUser(ctx, params)
}

// Authenticate checks credentials using sqlc GetUserByEmail
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

// AddRating inserts a rating via sqlc
func (s *Service) AddRating(ctx context.Context, userID sql.NullInt32, recipeID, rating int) (db.Rating, error) {
	params := db.InsertRatingParams{UserID: userID, RecipeID: sql.NullInt32{Int32: int32(recipeID), Valid: true}, Rating: sql.NullInt32{Int32: int32(rating), Valid: true}}
	return s.q.InsertRating(ctx, params)
}

// Favorites
func (s *Service) AddFavorite(ctx context.Context, userID int, recipeID int) (db.Favorite, error) {
	params := db.AddFavoriteParams{UserID: sql.NullInt32{Int32: int32(userID), Valid: true}, RecipeID: sql.NullInt32{Int32: int32(recipeID), Valid: true}}
	return s.q.AddFavorite(ctx, params)
}

func (s *Service) RemoveFavorite(ctx context.Context, userID int, recipeID int) error {
	params := db.RemoveFavoriteParams{UserID: sql.NullInt32{Int32: int32(userID), Valid: true}, RecipeID: sql.NullInt32{Int32: int32(recipeID), Valid: true}}
	return s.q.RemoveFavorite(ctx, params)
}

func (s *Service) ListFavorites(ctx context.Context, userID int) ([]db.ListFavoritesByUserRow, error) {
	return s.q.ListFavoritesByUser(ctx, sql.NullInt32{Int32: int32(userID), Valid: true})
}

// Suggestions: simple content-based filtering using tag overlap with user's favorites and high ratings
func (s *Service) GetSuggestions(ctx context.Context, userID int, limit int) ([]RecipeSummary, error) {
	favs, err := s.ListFavorites(ctx, userID)
	if err != nil {
		return nil, err
	}
	favoriteTagCounts := map[string]int{}
	for _, f := range favs {
		full, err := s.q.GetRecipeByID(ctx, f.RecipeID.Int32)
		if err != nil {
			// ignore missing
			continue
		}
		for _, t := range full.Tags {
			favoriteTagCounts[strings.ToLower(t)]++
		}
	}
	// fetch a broad set of recipes
	candidates, err := s.SearchAndFilterRecipes(ctx, "", "", "", nil, "", int(math.Max(float64(limit*5), 100)), 0)
	if err != nil {
		return nil, err
	}
	var scored []RecipeSummary
	for _, c := range candidates {
		score := 0
		for _, t := range c.Tags {
			score += favoriteTagCounts[strings.ToLower(t)]
		}
		if score > 0 {
			scored = append(scored, RecipeSummary{ID: c.ID, Title: c.Title, Score: score})
		}
	}
	sort.Slice(scored, func(i, j int) bool { return scored[i].Score > scored[j].Score })
	if len(scored) > limit {
		scored = scored[:limit]
	}
	return scored, nil
}

// Sentinel HTTP errors (optionally useful for handlers)
var (
	ErrBadRequest = fmt.Errorf("%d", http.StatusBadRequest)
)
