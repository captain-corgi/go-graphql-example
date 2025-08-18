-- Add authentication fields to users table
ALTER TABLE users 
ADD COLUMN password_hash VARCHAR(255) NOT NULL DEFAULT '',
ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT true,
ADD COLUMN last_login_at TIMESTAMP WITH TIME ZONE;

-- Create sessions table for refresh token management
CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_revoked BOOLEAN NOT NULL DEFAULT false,
    device_info TEXT,
    ip_address INET
);

-- Create indexes for sessions table
CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_refresh_token_hash ON user_sessions(refresh_token_hash);
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at);
CREATE INDEX idx_user_sessions_is_revoked ON user_sessions(is_revoked);

-- Create trigger for sessions updated_at
CREATE TRIGGER update_user_sessions_updated_at 
    BEFORE UPDATE ON user_sessions 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Clean up expired sessions function
CREATE OR REPLACE FUNCTION cleanup_expired_sessions()
RETURNS void AS $$
BEGIN
    DELETE FROM user_sessions 
    WHERE expires_at < NOW() OR is_revoked = true;
END;
$$ LANGUAGE plpgsql;

