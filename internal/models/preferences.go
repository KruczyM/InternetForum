package models

import "time"

type UserCategoryPreference struct {
	UserID     string
	CategoryID int
	Weight     int
}

type UserBookInteraction struct {
	UserID string
	BookID int
	Clicks int
	LastViewedAt time.Time
}