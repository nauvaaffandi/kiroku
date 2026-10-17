-- Migration 009: Shared Access
CREATE TABLE IF NOT EXISTS shares (
    id TEXT PRIMARY KEY,
    account_id TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_accessed_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_shares_account_id ON shares(account_id);
CREATE INDEX IF NOT EXISTS idx_shares_token ON shares(token);