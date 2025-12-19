package main

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"forum/internal/models"
)

func seedUsers(db *sql.DB) {
	users := []struct {
		Username  string
		FirstName string
		LastName  string
		Email     string
		Password  string
	}{
		{
			Username:  "admin",
			FirstName: "System",
			LastName:  "Administrator",
			Email:     "admin@literarylions.dev",
			Password:  "Admin123!",
		},
		{
			Username:  "alice",
			FirstName: "Alice",
			LastName:  "Reader",
			Email:     "alice.reader@mail.com",
			Password:  "Alice123!",
		},
		{
			Username:  "bob",
			FirstName: "Bob",
			LastName:  "Reader",
			Email:     "bob.reader@mail.com",
			Password:  "Bob123!",
		},
	}

	for _, u := range users {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("hash error: %v", err)
		}

		user := &models.User{
			ID:           uuid.NewString(),
			Username:     u.Username,
			FirstName:    u.FirstName,
			LastName:     u.LastName,
			Email:        u.Email,
			PasswordHash: string(hash),
		}

		if err := models.InsertUser(db, user); err != nil {
			log.Printf("cannot insert user %s: %v", u.Username, err)
		} else {
			log.Printf("inserted user: %s", u.Username)
		}
	}
}