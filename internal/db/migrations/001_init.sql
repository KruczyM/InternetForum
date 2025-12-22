PRAGMA foreign_keys = ON;

-- =========================
-- USERS & AUTH
-- =========================

CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,              -- UUID
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sessions (
    token TEXT PRIMARY KEY,
    data BLOB NOT NULL,
    expiry DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions(expiry);

-- =========================
-- BOOKS & CATEGORIES
-- =========================

CREATE TABLE IF NOT EXISTS books (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    kind TEXT NOT NULL                -- genre, theme, format, character, author
);

CREATE TABLE IF NOT EXISTS book_categories (
    book_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    PRIMARY KEY (book_id, category_id),
    FOREIGN KEY (book_id)
        REFERENCES books(id)
        ON DELETE CASCADE,
    FOREIGN KEY (category_id)
        REFERENCES categories(id)
        ON DELETE CASCADE
);

-- =========================
-- FORUM
-- =========================

CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    post_type TEXT NOT NULL DEFAULT 'discussion', -- discussion | analysis | review
    book_id INTEGER,                              -- nullable
    chapter TEXT,                                 -- nullable
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    FOREIGN KEY (book_id)
        REFERENCES books(id)
        ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS post_categories (
    post_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, category_id),
    FOREIGN KEY (post_id)
        REFERENCES posts(id)
        ON DELETE CASCADE,
    FOREIGN KEY (category_id)
        REFERENCES categories(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id)
        REFERENCES posts(id)
        ON DELETE CASCADE,
    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS likes (
    user_id TEXT NOT NULL,
    target_type TEXT NOT NULL,
    target_id INTEGER NOT NULL,
    value INTEGER NOT NULL CHECK (value IN (1, -1)),
    PRIMARY KEY (user_id, target_type, target_id),
    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- =========================
-- PREFERENCES & PERSONALIZATION
-- =========================

-- Explicit preferences (many-to-many, weighted)
CREATE TABLE IF NOT EXISTS user_category_preferences (
    user_id TEXT NOT NULL,
    category_id INTEGER NOT NULL,
    weight INTEGER NOT NULL DEFAULT 1 CHECK (weight BETWEEN 1 AND 5),
    PRIMARY KEY (user_id, category_id),
    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    FOREIGN KEY (category_id)
        REFERENCES categories(id)
        ON DELETE CASCADE
);

-- Implicit preferences (behavior-based)
CREATE TABLE IF NOT EXISTS user_book_interactions (
    user_id TEXT NOT NULL,
    book_id INTEGER NOT NULL,
    clicks INTEGER NOT NULL DEFAULT 1,
    last_viewed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, book_id),
    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    FOREIGN KEY (book_id)
        REFERENCES books(id)
        ON DELETE CASCADE
);

-- Optional ratings / likes
CREATE TABLE IF NOT EXISTS user_book_preferences (
    user_id TEXT NOT NULL,
    book_id INTEGER NOT NULL,
    rating INTEGER CHECK (rating BETWEEN 1 AND 5),
    liked BOOLEAN,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, book_id),
    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    FOREIGN KEY (book_id)
        REFERENCES books(id)
        ON DELETE CASCADE
);

-- =========================
-- CHAT
-- =========================

CREATE TABLE IF NOT EXISTS chat_messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);


INSERT OR IGNORE INTO users (id, email, username, first_name, last_name, password_hash) 
VALUES ('1', 'test@example.com', 'TestUser', 'Test', 'User', 'placeholder_hash');