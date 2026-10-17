-- Migration 002: Categories (hierarchical)
CREATE TABLE IF NOT EXISTS categories (
    id TEXT PRIMARY KEY,
    parent_id TEXT REFERENCES categories(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    icon TEXT,
    color TEXT,
    kind TEXT CHECK(kind IN ('must','have','want')),
    sort_order INTEGER NOT NULL DEFAULT 0,
    budget_amount INTEGER,
    budget_period TEXT CHECK(budget_period IN ('weekly','monthly','yearly')),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_categories_parent ON categories(parent_id);

CREATE TRIGGER IF NOT EXISTS categories_updated_at
AFTER UPDATE ON categories
BEGIN
    UPDATE categories SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- Prevent hard delete if transactions reference this category
CREATE TRIGGER IF NOT EXISTS categories_prevent_delete
BEFORE DELETE ON categories
WHEN EXISTS (SELECT 1 FROM transactions WHERE category_id = OLD.id AND deleted_at IS NULL)
BEGIN
    SELECT RAISE(ABORT, 'Category has associated transactions. Use soft delete (archived) instead.');
END;