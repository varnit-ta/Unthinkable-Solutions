-- name: AddFavorite :one
INSERT INTO favorites (user_id, recipe_id)
VALUES ($1, $2)
RETURNING id, user_id, recipe_id, created_at;

-- name: RemoveFavorite :exec
DELETE FROM favorites WHERE user_id = $1 AND recipe_id = $2;

-- name: ListFavoritesByUser :many
SELECT f.id as favorite_id, f.user_id, f.recipe_id, f.created_at, 
  r.title, r.description, r.cuisine, r.difficulty, r.diet_type, 
  r.prep_time_minutes, r.cook_time_minutes, r.total_time_minutes, r.servings,
  COALESCE((SELECT ROUND(AVG(rating)::numeric, 1) FROM ratings WHERE recipe_id = r.id)::text, '0') as average_rating
FROM favorites f
JOIN recipes r ON r.id = f.recipe_id
WHERE f.user_id = $1
ORDER BY f.created_at DESC;

-- name: IsFavorite :one
SELECT EXISTS(
  SELECT 1 FROM favorites
  WHERE user_id = $1 AND recipe_id = $2
) as is_favorite;
