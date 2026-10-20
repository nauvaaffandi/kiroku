-- Migration 014: Add archived column to categories
ALTER TABLE categories ADD COLUMN archived INTEGER NOT NULL DEFAULT 0;
