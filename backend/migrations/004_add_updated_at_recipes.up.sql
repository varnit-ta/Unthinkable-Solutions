-- Add updated_at column to recipes
ALTER TABLE recipes
  ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT now();
