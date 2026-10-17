-- Migration 010: Telegram Bindings
CREATE TABLE IF NOT EXISTS telegram_bindings (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    telegram_chat_id TEXT NOT NULL UNIQUE,
    telegram_user_id TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
