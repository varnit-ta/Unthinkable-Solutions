-- Rollback users and favorites changes
DROP TABLE IF EXISTS favorites CASCADE;

ALTER TABLE users
  DROP COLUMN IF EXISTS password_hash,
  DROP COLUMN IF EXISTS email;
