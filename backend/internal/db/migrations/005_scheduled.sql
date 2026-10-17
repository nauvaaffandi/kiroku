-- Migration 005: Scheduled/Recurring Transactions
CREATE TABLE IF NOT EXISTS scheduled (
    id TEXT PRIMARY KEY,
    rrule TEXT NOT NULL,
    next_run DATETIME NOT NULL,
    max_occurrences INTEGER,
    occurrences_done INTEGER NOT NULL DEFAULT 0,
    transaction_template_json TEXT NOT NULL,
    enabled INTEGER NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_scheduled_next_run ON scheduled(next_run) WHERE enabled = 1;
CREATE INDEX IF NOT EXISTS idx_scheduled_enabled ON scheduled(enabled);