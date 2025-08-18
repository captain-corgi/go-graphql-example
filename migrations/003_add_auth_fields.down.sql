-- Drop cleanup function
DROP FUNCTION IF EXISTS cleanup_expired_sessions();

-- Drop sessions table
DROP TABLE IF EXISTS user_sessions;

-- Remove authentication fields from users table
ALTER TABLE users 
DROP COLUMN IF EXISTS password_hash,
DROP COLUMN IF EXISTS is_active,
DROP COLUMN IF EXISTS last_login_at;

