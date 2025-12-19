package models

import "time"

type ChatMessage struct {
	ID        int
	UserID    string
	Content   string
	CreatedAt time.Time
}