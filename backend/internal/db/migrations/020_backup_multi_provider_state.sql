-- 020_backup_multi_provider_state.sql
CREATE TABLE IF NOT EXISTS backup_provider_sync_states (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    execution_id TEXT NOT NULL,
    provider_name TEXT NOT NULL,
    status TEXT NOT NULL, -- 'success', 'failed'
    error_message TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_backup_provider_sync_exec ON backup_provider_sync_states(execution_id);
