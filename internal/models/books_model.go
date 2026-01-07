package models

import (
	"database/sql"
	"time"
)

type Book struct {
	ID          int
	Title       string
	Author      string
	Description string
	CreatedAt   time.Time
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

func (m *BookModel) GetAllBooks() ([]Book, error) {
	rows, err := m.DB.Query(`
		SELECT id, title
		FROM books
		ORDER BY title ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Title); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, rows.Err()
}
