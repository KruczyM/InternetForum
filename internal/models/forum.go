package models

import "time"

type Post struct {
	ID int
	UserID string
	Title   string
	Content string
	PostType string // "discussion", "analysis", "review"
	BookID  *int    // can be null
	Chapter *string // can be null
	CreatedAt time.Time
}
type Comment struct {
	ID        int
	PostID    int
	UserID    string
	Content   string
	CreatedAt time.Time
}

// one vote per person per post + 1 or -1
type Vote struct {
	ID         int
	UserID     string
	TargetType string // "post" | "comment"
	TargetID   int
	Value      int    // +1 or -1
}