-- Migration 021: Remove budget_amount & budget_period from categories
-- The goals/budget system has been refactored out. Categories are now pure.
ALTER TABLE categories DROP COLUMN budget_amount;
ALTER TABLE categories DROP COLUMN budget_period;
