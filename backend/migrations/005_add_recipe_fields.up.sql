-- Add missing recipe fields
ALTER TABLE recipes
ADD COLUMN IF NOT EXISTS description TEXT,
ADD COLUMN IF NOT EXISTS diet_type TEXT,
ADD COLUMN IF NOT EXISTS prep_time_minutes INTEGER,
ADD COLUMN IF NOT EXISTS total_time_minutes INTEGER;

-- Update total_time_minutes to be sum of prep and cook time where possible
UPDATE recipes
SET total_time_minutes = COALESCE(prep_time_minutes, 0) + COALESCE(cook_time_minutes, 0)
WHERE total_time_minutes IS NULL;
