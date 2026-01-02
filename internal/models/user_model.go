package models

import (
	"database/sql"
	"strings"
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
	if err != nil {
		// SQLite UNIQUE constraint
		if strings.Contains(err.Error(), "users.email") {
			return ErrDuplicateEmail
		}
		if strings.Contains(err.Error(), "users.username") {
			return ErrDuplicateUsername
		}
		return err
	}

	return nil
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
	if err := row.Scan(&user.ID, &user.PasswordHash); err != nil {
		return nil, err
	}
	return &user, nil
}
// retrieves a user by id
func GetUserByID(db *sql.DB, id string) (*User, error) {
	query := `
	SELECT id, email,avatar_path, username, first_name, last_name, password_hash, created_at
	FROM users 
	WHERE id = ?`
	var user User
	if err := db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.AvatarPath, &user.Username, &user.FirstName, &user.LastName,&user.PasswordHash, &user.CreatedAt,); err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUserNameFields(db *sql.DB, userID, firstName, lastName string) error {
	query := `UPDATE users SET first_name = ?, last_name = ? WHERE id = ?`
	_, err := db.Exec(query, firstName, lastName, userID)
	return err
}

func UpdateUserPasswordHash(db *sql.DB, userID, newHash string) error {
	query := `UPDATE users SET password_hash = ? WHERE id = ?`
	_, err := db.Exec(query, newHash, userID)
	return err
}

func UpdateUserAvatarPath(db *sql.DB, userID, newPath string) error {
	query := `UPDATE users SET avatar_path = ? WHERE id = ?`
	_, err := db.Exec(query, newPath, userID)
	return err
}

func CountUserLikes(db *sql.DB, userID string) (int, error) {
	query := `SELECT COUNT(*) FROM likes WHERE user_id = ?`
	var count int
	err := db.QueryRow(query, userID).Scan(&count)
	return count, err
}

func CountUserPosts(db *sql.DB, userID string) (int, error) {
	query := `SELECT COUNT(*) FROM posts WHERE user_id = ?`
	var count int
	err := db.QueryRow(query, userID).Scan(&count)
	return count, err
}

func CountUserComments(db *sql.DB, userID string) (int, error) {
	query := `SELECT COUNT(*) FROM comments WHERE user_id = ?`
	var count int
	err := db.QueryRow(query, userID).Scan(&count)
	return count, err
}


func ExistsUser(db *sql.DB, id string) (bool, error){
	var exists bool
	stmt :=`
	SELECT EXISTS
	(SELECT true 
	FROM users 
	WHERE id = ?)`
	err := db.QueryRow(stmt,id).Scan(&exists)
	return exists, err
}




