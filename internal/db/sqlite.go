package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return nil, err
	}

	return db, nil
}

func InitDB(path string) (*sql.DB, error) {
	db, err := Open(path)
	if err != nil {
		return nil, err
	}

	if err := RunMigrations(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}