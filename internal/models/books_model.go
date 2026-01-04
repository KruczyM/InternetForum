package models

import (
	"time"
	"database/sql"
)

type Book struct {
	ID          int
	Title       string
	Author      string
	Description string
	CreatedAt time.Time
}

// BookCategory is a many-to-many relationship between Book and Category
type BookCategory struct {
	BookID     int
	CategoryID int
}

type Category struct {
	ID   int
	Name string
	Kind string // "genre", "theme", "format", "character", "author"
}

type BookModel struct {
	DB *sql.DB
}

func (m *BookModel) AddBook(title, author, description string) error {
	stmt := `INSERT INTO books (title, author, description, created_at)
    VALUES(?, ?, ?, datetime('now'))`

	_, err := m.DB.Exec(stmt, title, author, description)
	return err
}