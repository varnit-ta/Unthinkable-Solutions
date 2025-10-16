-- name: ListRecipes :many
-- List a page of recipes
SELECT id, title, description, cuisine, difficulty, diet_type, prep_time_minutes, cook_time_minutes, total_time_minutes, servings,
  COALESCE((SELECT ROUND(AVG(rating)::numeric, 1)::text FROM ratings r WHERE r.recipe_id = recipes.id), '0') as average_rating
FROM recipes
ORDER BY recipes.id
LIMIT $1 OFFSET $2;

-- name: GetRecipeByID :one
SELECT id, title, description, cuisine, difficulty, diet_type, prep_time_minutes, cook_time_minutes, total_time_minutes, servings, ingredients, steps, nutrition, tags,
  COALESCE((SELECT ROUND(AVG(rating)::numeric, 1)::text FROM ratings r WHERE r.recipe_id = recipes.id), '0') as average_rating
FROM recipes
WHERE recipes.id = $1;

-- name: InsertRating :one
INSERT INTO ratings (user_id, recipe_id, rating)
VALUES ($1, $2, $3)
RETURNING id, user_id, recipe_id, rating, created_at;

-- name: GetRatingsForRecipe :many
SELECT id, user_id, recipe_id, rating, created_at
FROM ratings
WHERE ratings.recipe_id = $1;

-- name: SearchRecipes :many
-- Simple search by title or tags
SELECT id, title, description, cuisine, difficulty, diet_type, prep_time_minutes, cook_time_minutes, total_time_minutes, servings, ingredients, steps, nutrition, tags,
  COALESCE((SELECT ROUND(AVG(rating)::numeric, 1)::text FROM ratings r WHERE r.recipe_id = recipes.id), '0') as average_rating
FROM recipes
WHERE recipes.title ILIKE '%' || $1 || '%' OR $1 = ANY(recipes.tags)
ORDER BY recipes.id
LIMIT $2 OFFSET $3;

-- name: CreateRecipe :one
INSERT INTO recipes (title, description, cuisine, difficulty, diet_type, prep_time_minutes, cook_time_minutes, total_time_minutes, servings, tags, ingredients, steps, nutrition)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING id;
