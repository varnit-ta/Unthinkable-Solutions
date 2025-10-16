package handlers

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/sqlc-dev/pqtype"
	"github.com/varnit-ta/smart-recipe-generator/backend/internal/db"
)

// RecipeListResponse is a clean JSON response for recipe lists
type RecipeListResponse struct {
	ID               int32    `json:"id"`
	Title            string   `json:"title"`
	Description      string   `json:"description,omitempty"`
	Cuisine          string   `json:"cuisine,omitempty"`
	Difficulty       string   `json:"difficulty,omitempty"`
	DietType         string   `json:"diet_type,omitempty"`
	PrepTimeMinutes  int      `json:"prep_time_minutes,omitempty"`
	CookTimeMinutes  int      `json:"cook_time_minutes,omitempty"`
	TotalTimeMinutes int      `json:"total_time_minutes,omitempty"`
	Servings         int      `json:"servings,omitempty"`
	AverageRating    string   `json:"average_rating"`
	Tags             []string `json:"tags,omitempty"`
}

// RecipeDetailResponse is a clean JSON response for full recipe details
type RecipeDetailResponse struct {
	ID               int32       `json:"id"`
	Title            string      `json:"title"`
	Description      string      `json:"description,omitempty"`
	Cuisine          string      `json:"cuisine,omitempty"`
	Difficulty       string      `json:"difficulty,omitempty"`
	DietType         string      `json:"diet_type,omitempty"`
	PrepTimeMinutes  int         `json:"prep_time_minutes,omitempty"`
	CookTimeMinutes  int         `json:"cook_time_minutes,omitempty"`
	TotalTimeMinutes int         `json:"total_time_minutes,omitempty"`
	Servings         int         `json:"servings,omitempty"`
	Ingredients      interface{} `json:"ingredients,omitempty"`
	Steps            interface{} `json:"steps,omitempty"`
	Nutrition        interface{} `json:"nutrition,omitempty"`
	Tags             []string    `json:"tags,omitempty"`
	AverageRating    string      `json:"average_rating"`
}

func toRecipeListResponse(row db.ListRecipesRow) RecipeListResponse {
	return RecipeListResponse{
		ID:               row.ID,
		Title:            row.Title,
		Description:      nullStringValue(row.Description),
		Cuisine:          nullStringValue(row.Cuisine),
		Difficulty:       nullStringValue(row.Difficulty),
		DietType:         nullStringValue(row.DietType),
		PrepTimeMinutes:  int(nullInt32Value(row.PrepTimeMinutes)),
		CookTimeMinutes:  int(nullInt32Value(row.CookTimeMinutes)),
		TotalTimeMinutes: int(nullInt32Value(row.TotalTimeMinutes)),
		Servings:         int(nullInt32Value(row.Servings)),
		AverageRating:    interfaceToString(row.AverageRating),
	}
}

func toRecipeDetailResponse(row db.GetRecipeByIDRow) RecipeDetailResponse {
	return RecipeDetailResponse{
		ID:               row.ID,
		Title:            row.Title,
		Description:      nullStringValue(row.Description),
		Cuisine:          nullStringValue(row.Cuisine),
		Difficulty:       nullStringValue(row.Difficulty),
		DietType:         nullStringValue(row.DietType),
		PrepTimeMinutes:  int(nullInt32Value(row.PrepTimeMinutes)),
		CookTimeMinutes:  int(nullInt32Value(row.CookTimeMinutes)),
		TotalTimeMinutes: int(nullInt32Value(row.TotalTimeMinutes)),
		Servings:         int(nullInt32Value(row.Servings)),
		Ingredients:      pqNullRawMessageValue(row.Ingredients),
		Steps:            pqNullRawMessageValue(row.Steps),
		Nutrition:        pqNullRawMessageValue(row.Nutrition),
		Tags:             row.Tags,
		AverageRating:    interfaceToString(row.AverageRating),
	}
}

func toSearchRecipeResponse(row db.SearchRecipesRow) RecipeDetailResponse {
	return RecipeDetailResponse{
		ID:               row.ID,
		Title:            row.Title,
		Description:      nullStringValue(row.Description),
		Cuisine:          nullStringValue(row.Cuisine),
		Difficulty:       nullStringValue(row.Difficulty),
		DietType:         nullStringValue(row.DietType),
		PrepTimeMinutes:  int(nullInt32Value(row.PrepTimeMinutes)),
		CookTimeMinutes:  int(nullInt32Value(row.CookTimeMinutes)),
		TotalTimeMinutes: int(nullInt32Value(row.TotalTimeMinutes)),
		Servings:         int(nullInt32Value(row.Servings)),
		Ingredients:      pqNullRawMessageValue(row.Ingredients),
		Steps:            pqNullRawMessageValue(row.Steps),
		Nutrition:        pqNullRawMessageValue(row.Nutrition),
		Tags:             row.Tags,
		AverageRating:    interfaceToString(row.AverageRating),
	}
}

func nullStringValue(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func nullInt32Value(ni sql.NullInt32) int32 {
	if ni.Valid {
		return ni.Int32
	}
	return 0
}

func pqNullRawMessageValue(nrm pqtype.NullRawMessage) interface{} {
	if nrm.Valid {
		return nrm.RawMessage
	}
	return nil
}

func interfaceToString(i interface{}) string {
	if i == nil {
		return "0"
	}
	if s, ok := i.(string); ok {
		return s
	}
	return "0"
}

func nullStringArrayValue(arr []sql.NullString) []string {
	result := make([]string, 0, len(arr))
	for _, ns := range arr {
		if ns.Valid {
			result = append(result, ns.String)
		}
	}
	return result
}

func pqStringArrayValue(arr pq.StringArray) []string {
	if arr == nil {
		return []string{}
	}
	return []string(arr)
}

type FavoriteRecipeResponse struct {
	FavoriteID       int32  `json:"favorite_id"`
	RecipeID         int32  `json:"recipe_id"`
	Title            string `json:"title"`
	Description      string `json:"description,omitempty"`
	Cuisine          string `json:"cuisine,omitempty"`
	Difficulty       string `json:"difficulty,omitempty"`
	DietType         string `json:"diet_type,omitempty"`
	PrepTimeMinutes  int    `json:"prep_time_minutes,omitempty"`
	CookTimeMinutes  int    `json:"cook_time_minutes,omitempty"`
	TotalTimeMinutes int    `json:"total_time_minutes,omitempty"`
	Servings         int    `json:"servings,omitempty"`
	AverageRating    string `json:"average_rating"`
}

func toFavoriteRecipeResponse(row db.ListFavoritesByUserRow) FavoriteRecipeResponse {
	return FavoriteRecipeResponse{
		FavoriteID:       row.FavoriteID,
		RecipeID:         nullInt32Value(row.RecipeID),
		Title:            row.Title,
		Description:      nullStringValue(row.Description),
		Cuisine:          nullStringValue(row.Cuisine),
		Difficulty:       nullStringValue(row.Difficulty),
		DietType:         nullStringValue(row.DietType),
		PrepTimeMinutes:  int(nullInt32Value(row.PrepTimeMinutes)),
		CookTimeMinutes:  int(nullInt32Value(row.CookTimeMinutes)),
		TotalTimeMinutes: int(nullInt32Value(row.TotalTimeMinutes)),
		Servings:         int(nullInt32Value(row.Servings)),
		AverageRating:    row.AverageRating,
	}
}
