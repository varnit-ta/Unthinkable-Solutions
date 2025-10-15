-- name: ListRecipes :many
-- List a page of recipes
SELECT id, title, cuisine, difficulty, cook_time_minutes, servings
FROM recipes
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: GetRecipeByID :one
SELECT id, title, cuisine, difficulty, cook_time_minutes, servings, ingredients, steps, nutrition, tags
FROM recipes
WHERE id = $1;

-- name: InsertRating :one
INSERT INTO ratings (user_id, recipe_id, rating)
VALUES ($1, $2, $3)
RETURNING id, user_id, recipe_id, rating, created_at;

-- name: GetRatingsForRecipe :many
SELECT id, user_id, recipe_id, rating, created_at
FROM ratings
WHERE recipe_id = $1;

-- name: SearchRecipes :many
-- Simple search by title or tags
SELECT id, title, cuisine, difficulty, cook_time_minutes, servings, ingredients, steps, nutrition, tags
FROM recipes
WHERE title ILIKE '%' || $1 || '%' OR $1 = ANY(tags)
ORDER BY id
LIMIT $2 OFFSET $3;

-- name: CreateRecipe :one
INSERT INTO recipes (title, cuisine, difficulty, cook_time_minutes, servings, tags, ingredients, steps, nutrition)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id;
