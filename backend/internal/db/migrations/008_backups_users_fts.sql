-- Migration 008: Backups & Users & FTS5
CREATE TABLE IF NOT EXISTS backups (
    id TEXT PRIMARY KEY,
    filename TEXT NOT NULL,
    size INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    retention_tier TEXT NOT NULL CHECK(retention_tier IN ('daily','weekly','monthly'))
);

CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    totp_secret TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- FTS5 for transaction search
CREATE VIRTUAL TABLE IF NOT EXISTS transactions_fts USING fts5(
    note, payee_payer, labels,
    content='transactions',
    content_rowid='rowid'
);

CREATE TRIGGER IF NOT EXISTS transactions_ai AFTER INSERT ON transactions BEGIN
    INSERT INTO transactions_fts(rowid, note, payee_payer, labels) VALUES (new.rowid, new.note, new.payee_payer, new.labels);
END;

CREATE TRIGGER IF NOT EXISTS transactions_ad AFTER DELETE ON transactions BEGIN
    INSERT INTO transactions_fts(transactions_fts, rowid, note, payee_payer, labels) VALUES ('delete', old.rowid, old.note, old.payee_payer, old.labels);
END;

CREATE TRIGGER IF NOT EXISTS transactions_au AFTER UPDATE ON transactions BEGIN
    INSERT INTO transactions_fts(transactions_fts, rowid, note, payee_payer, labels) VALUES ('delete', old.rowid, old.note, old.payee_payer, old.labels);
    INSERT INTO transactions_fts(rowid, note, payee_payer, labels) VALUES (new.rowid, new.note, new.payee_payer, new.labels);
END;