package models

import (
	"time"

)

type User struct {
	ID        string    
	Email     string
	Username  string
	FirstName string
	LastName  string

	PasswordHash string

	CreatedAt time.Time
}


type UserBookPreference struct {
	UserID string
	BookID int
	Rating *int
	Liked  *bool
	CreatedAt time.Time
}