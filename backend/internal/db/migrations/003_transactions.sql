-- Migration 003: Transactions
CREATE TABLE IF NOT EXISTS transactions (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL CHECK(type IN ('income','expense','transfer')),
    date_time DATETIME NOT NULL,
    amount INTEGER NOT NULL,
    currency TEXT NOT NULL CHECK(length(currency) = 3 AND currency GLOB '[A-Z][A-Z][A-Z]'),
    from_account_id TEXT REFERENCES accounts(id) ON DELETE SET NULL,
    to_account_id TEXT REFERENCES accounts(id) ON DELETE SET NULL,
    category_id TEXT REFERENCES categories(id) ON DELETE SET NULL,
    labels TEXT,
    note TEXT,
    payee_payer TEXT,
    payment_type TEXT CHECK(payment_type IN ('cash','card','bank','ewallet','crypto','other')),
    warranty_until DATE,
    status TEXT NOT NULL DEFAULT 'cleared' CHECK(status IN ('cleared','uncleared','reconciled')),
    place_name TEXT,
    place_lat REAL,
    place_lng REAL,
    attachments TEXT,
    template_id TEXT REFERENCES templates(id) ON DELETE SET NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_transactions_date_time ON transactions(date_time);
CREATE INDEX IF NOT EXISTS idx_transactions_from_account ON transactions(from_account_id);
CREATE INDEX IF NOT EXISTS idx_transactions_to_account ON transactions(to_account_id);
CREATE INDEX IF NOT EXISTS idx_transactions_category ON transactions(category_id);
CREATE INDEX IF NOT EXISTS idx_transactions_deleted ON transactions(deleted_at) WHERE deleted_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type);

CREATE TRIGGER IF NOT EXISTS transactions_updated_at
AFTER UPDATE ON transactions
BEGIN
    UPDATE transactions SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;