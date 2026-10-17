-- Migration 011: Password Change Tokens
CREATE TABLE IF NOT EXISTS password_change_tokens (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    token TEXT NOT NULL UNIQUE,
    type TEXT NOT NULL CHECK(type IN ('register', 'changepassword', 'forgotpassword')),
    expires_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_password_change_tokens_token ON password_change_tokens(token);
CREATE INDEX IF NOT EXISTS idx_password_change_tokens_expires ON password_change_tokens(expires_at);