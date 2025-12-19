package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
)


//hashes a password using SHA-256 (for demo purposes only)
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// compares a plain password with a hashed password
func CheckPassword(password, hashed string) bool {
	return HashPassword(password) == hashed
}

// inserts a new user into the database
func InsertUser(db *sql.DB, user *User) ( error ) {
	// id will be genrated by UUID package - so it will be random number
	query := `INSERT INTO users (id, username, first_name, last_name, email, password_hash) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, user.ID, user.Username, user.FirstName, user.LastName, user.Email, user.PasswordHash)
	return err
}

// retrieves a user by username
func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	query := `SELECT id, username, first_name, last_name, email, password_hash FROM users WHERE username = ?`
	row := db.QueryRow(query, username)
	var user User
	if err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash); err != nil {
		return nil, err
	}
	return &user, nil
}
