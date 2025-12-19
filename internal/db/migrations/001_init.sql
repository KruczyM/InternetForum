PRAGMA foreign_keys = ON;

-- =========================
-- USERS & AUTH
-- =========================

CREATE TABLE users (
    id TEXT PRIMARY KEY,              -- UUID
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
    id TEXT PRIMARY KEY,              -- UUID
    user_id TEXT NOT NULL,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- =========================
-- BOOKS & CATEGORIES
-- =========================

CREATE TABLE books (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    kind TEXT NOT NULL                -- genre, theme, format, character, author
);

CREATE TABLE book_categories (
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

CREATE TABLE posts (
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

CREATE TABLE post_categories (
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

CREATE TABLE comments (
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

CREATE TABLE votes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    target_type TEXT NOT NULL,         -- post | comment
    target_id INTEGER NOT NULL,
    value INTEGER NOT NULL CHECK (value IN (1, -1)),
    UNIQUE (user_id, target_type, target_id),
    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- =========================
-- PREFERENCES & PERSONALIZATION
-- =========================

-- Explicit preferences (many-to-many, weighted)
CREATE TABLE user_category_preferences (
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
CREATE TABLE user_book_interactions (
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
CREATE TABLE user_book_preferences (
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

CREATE TABLE chat_messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);
