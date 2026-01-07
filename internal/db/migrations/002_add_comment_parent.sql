PRAGMA foreign_keys = ON;

-- Add parent_id column to comments to support replies
ALTER TABLE comments ADD COLUMN parent_id INTEGER;
ALTER TABLE users ADD COLUMN about TEXT DEFAULT '';

-- Note: SQLite doesn't support adding foreign key constraints via ALTER TABLE easily.
-- The column is nullable and can store the parent comment id for threaded replies.
