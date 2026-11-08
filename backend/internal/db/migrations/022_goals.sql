-- Migration 022: Goals table
CREATE TABLE IF NOT EXISTS goals (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    target_amount INTEGER NOT NULL,
    icon TEXT,
    desired_date DATE,
    color_theme TEXT,
    note TEXT,
    saved_already INTEGER NOT NULL DEFAULT 0,
    currency TEXT NOT NULL,
    deleted_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_goals_deleted ON goals(deleted_at) WHERE deleted_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_goals_currency ON goals(currency);

CREATE TRIGGER IF NOT EXISTS goals_updated_at
AFTER UPDATE ON goals
BEGIN
    UPDATE goals SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
