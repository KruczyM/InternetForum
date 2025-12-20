package models

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
)


//hashed password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// compares a plain password with a hashed password
func CheckPassword(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashed),
		[]byte(password),
	)
	return err == nil
}

// inserts a new user into the database
func InsertUser(db *sql.DB, user *User) ( error ) {
	// id will be genrated by UUID package - so it will be random number
	query := `
	INSERT INTO users (id, username, first_name, last_name, email, password_hash) 
	VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, user.ID, user.Username, user.FirstName, user.LastName, user.Email, user.PasswordHash)
	return err
}

// retrieves a user by username
func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	query := `
	SELECT id, username, first_name, last_name, email, password_hash 
	FROM users 
	WHERE username = ?`
	row := db.QueryRow(query, username)
	var user User
	if err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash); err != nil {
		return nil, err
	}
	return &user, nil
}

// retrieves a user by email
func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	query := `
	SELECT id, password_hash 
	FROM users 
	WHERE email = ?`
	row := db.QueryRow(query, email)
	var user User
	if err := row.Scan(&user.Email, &user.PasswordHash); err != nil {
		return nil, err
	}
	return &user, nil
}
