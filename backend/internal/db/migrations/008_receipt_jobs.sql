-- Migration 008: Receipt Processing Jobs
CREATE TABLE IF NOT EXISTS receipt_jobs (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'created',
    stage TEXT NOT NULL DEFAULT 'init',
    original_size INTEGER,
    compressed_size INTEGER,
    ocr_text TEXT,
    ai_raw_output TEXT,
    transaction_draft TEXT,
    transaction_id TEXT REFERENCES transactions(id) ON DELETE SET NULL,
    error_message TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS receipt_job_events (
    id TEXT PRIMARY KEY,
    job_id TEXT NOT NULL REFERENCES receipt_jobs(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    stage TEXT NOT NULL,
    attempt INTEGER DEFAULT 1,
    status TEXT NOT NULL,
    message TEXT,
    metadata TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_receipt_jobs_user ON receipt_jobs(user_id);
CREATE INDEX IF NOT EXISTS idx_receipt_job_events_job ON receipt_job_events(job_id);

CREATE TRIGGER IF NOT EXISTS receipt_jobs_updated_at
AFTER UPDATE ON receipt_jobs
BEGIN
    UPDATE receipt_jobs SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;