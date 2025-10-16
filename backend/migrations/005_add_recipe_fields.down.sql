-- Rollback added recipe fields
ALTER TABLE recipes
DROP COLUMN IF EXISTS description,
DROP COLUMN IF EXISTS diet_type,
DROP COLUMN IF EXISTS prep_time_minutes,
DROP COLUMN IF EXISTS total_time_minutes;
