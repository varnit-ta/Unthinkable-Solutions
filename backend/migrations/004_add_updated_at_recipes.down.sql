-- Remove updated_at column from recipes
ALTER TABLE recipes
  DROP COLUMN IF EXISTS updated_at;
