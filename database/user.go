package database

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
)

var dbInstance *sql.DB

// opens the SQLite database and creates the users table
func InitDB(dataSourceName string) error {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	dbInstance = db
	return CreateUserTable(db)
}

// returns the global DB instance
func GetDB() *sql.DB {
	return dbInstance
}

//hashes a password using SHA-256 (for demo purposes only)
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// compares a plain password with a hashed password
func CheckPassword(password, hashed string) bool {
	return HashPassword(password) == hashed
}

type User struct {
	ID       int
	Username string
	Name     string
	Surname  string
	FullName string
	Email    string
	Password string // Store hashed password
}

// creates the users table if it doesn't exist
func CreateUserTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		surname TEXT NOT NULL,
		fullname TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`
	_, err := db.Exec(query)
	return err
}

// inserts a new user into the database
func InsertUser(db *sql.DB, user *User) (int64, error) {
	query := `INSERT INTO users (username, name, surname, fullname, email, password) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, user.Username, user.Name, user.Surname, user.FullName, user.Email, user.Password)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	user.ID = int(id)
	return id, nil
}

// retrieves a user by username
func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	query := `SELECT id, username, name, surname, fullname, email, password FROM users WHERE username = ?`
	row := db.QueryRow(query, username)
	var user User
	if err := row.Scan(&user.ID, &user.Username, &user.Name, &user.Surname, &user.FullName, &user.Email, &user.Password); err != nil {
		return nil, err
	}
	return &user, nil
}
