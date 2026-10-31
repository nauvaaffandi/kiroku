-- 019_backup_provider_state.sql
CREATE TABLE IF NOT EXISTS backup_provider_states (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    execution_id TEXT NOT NULL,
    provider_name TEXT NOT NULL,
    artifact_key TEXT NOT NULL,
    status TEXT NOT NULL,
    error_message TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_backup_states_exec ON backup_provider_states(execution_id);
