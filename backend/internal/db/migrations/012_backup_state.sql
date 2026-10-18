-- Migration 012: Backup scheduling state (XLSX backup system)
CREATE TABLE IF NOT EXISTS backup_state (
    id TEXT PRIMARY KEY,
    last_created_processed DATE,
    last_run_at DATETIME,
    last_status TEXT,
    last_error TEXT,
    last_pushed_at DATETIME,
    last_commit TEXT,
    format_version INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS backup_dates (
    date TEXT PRIMARY KEY,
    fingerprint TEXT NOT NULL,
    has_data INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'processed',
    checked_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS backup_months (
    month TEXT PRIMARY KEY,
    fingerprint TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'processed',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS backup_rechecks (
    date TEXT PRIMARY KEY,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO backup_state (id, last_created_processed, last_run_at, last_status, last_error, last_pushed_at, last_commit, format_version)
VALUES ('singleton', NULL, NULL, NULL, NULL, NULL, NULL, 1)
ON CONFLICT(id) DO NOTHING;
