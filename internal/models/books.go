package models

import "time"

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

