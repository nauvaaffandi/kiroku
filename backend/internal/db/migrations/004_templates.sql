-- Migration 004: Templates
CREATE TABLE IF NOT EXISTS templates (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL CHECK(type IN ('income','expense','transfer')),
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
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);