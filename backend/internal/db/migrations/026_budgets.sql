-- Migration 026: Budgets (transaction-driven spending limits)
CREATE TABLE IF NOT EXISTS budgets (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    period TEXT NOT NULL CHECK(period IN ('weekly','monthly','yearly','one_time')),
    amount INTEGER NOT NULL CHECK(amount > 0),
    currency TEXT NOT NULL CHECK(length(currency) = 3 AND currency GLOB '[A-Z][A-Z][A-Z]'),
    account_id TEXT REFERENCES accounts(id) ON DELETE SET NULL,
    trending_over INTEGER NOT NULL DEFAULT 0,
    over_budget INTEGER NOT NULL DEFAULT 0,
    starts_at DATETIME,
    ends_at DATETIME,
    deleted_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_budgets_user ON budgets(user_id);
CREATE INDEX IF NOT EXISTS idx_budgets_user_period ON budgets(user_id, period);
CREATE INDEX IF NOT EXISTS idx_budgets_account ON budgets(account_id);
CREATE INDEX IF NOT EXISTS idx_budgets_deleted ON budgets(deleted_at) WHERE deleted_at IS NOT NULL;

CREATE TABLE IF NOT EXISTS budget_categories (
    budget_id TEXT NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    category_id TEXT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (budget_id, category_id)
);

CREATE INDEX IF NOT EXISTS idx_budget_categories_category ON budget_categories(category_id);

CREATE TABLE IF NOT EXISTS budget_labels (
    budget_id TEXT NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    label_id TEXT NOT NULL REFERENCES labels(id) ON DELETE CASCADE,
    PRIMARY KEY (budget_id, label_id)
);

CREATE INDEX IF NOT EXISTS idx_budget_labels_label ON budget_labels(label_id);

CREATE TRIGGER IF NOT EXISTS budgets_updated_at
AFTER UPDATE ON budgets
BEGIN
    UPDATE budgets SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- Prevent duplicate category and label associations are enforced by PRIMARY KEY.
