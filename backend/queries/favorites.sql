-- name: AddFavorite :one
INSERT INTO favorites (user_id, recipe_id)
VALUES ($1, $2)
RETURNING id, user_id, recipe_id, created_at;

-- name: RemoveFavorite :exec
DELETE FROM favorites WHERE user_id = $1 AND recipe_id = $2;

-- name: ListFavoritesByUser :many
SELECT f.id, f.user_id, f.recipe_id, f.created_at, r.title, r.cuisine
FROM favorites f
JOIN recipes r ON r.id = f.recipe_id
WHERE f.user_id = $1
ORDER BY f.created_at DESC;
