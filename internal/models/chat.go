package models

import "time"

type ChatMessage struct {
	ID        int
	UserID    string
	Content   string
	Category  string // Added for filtering by category
	CreatedAt time.Time
}