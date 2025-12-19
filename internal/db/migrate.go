package database

import (
	"database/sql"
	"os"
)

func RunMigrations(db *sql.DB) error {
	files := []string{
		"internal/db/migrations/001_init.sql",
	}

	for _, file := range files {
		query, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		if _, err := db.Exec(string(query)); err != nil {
			return err
		}
	}

	return nil
}
