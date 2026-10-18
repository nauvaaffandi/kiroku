-- Migration 013: Rename password -> PIN
ALTER TABLE users RENAME COLUMN password_hash TO pin_hash;

ALTER TABLE password_change_tokens RENAME TO pin_change_tokens;

DROP INDEX IF EXISTS idx_password_change_tokens_token;
DROP INDEX IF EXISTS idx_password_change_tokens_expires;
CREATE INDEX IF NOT EXISTS idx_pin_change_tokens_token ON pin_change_tokens(token);
CREATE INDEX IF NOT EXISTS idx_pin_change_tokens_expires ON pin_change_tokens(expires_at);
