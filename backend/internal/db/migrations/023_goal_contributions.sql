-- Migration 023: Goal contributions table
CREATE TABLE IF NOT EXISTS goal_contributions (
    id TEXT PRIMARY KEY,
    goal_id TEXT NOT NULL REFERENCES goals(id) ON DELETE CASCADE,
    amount INTEGER NOT NULL,
    added_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    note TEXT
);

CREATE INDEX IF NOT EXISTS idx_goal_contributions_goal_id ON goal_contributions(goal_id);

CREATE TRIGGER IF NOT EXISTS goal_contributions_updated_at
AFTER UPDATE ON goal_contributions
BEGIN
    UPDATE goal_contributions SET added_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
