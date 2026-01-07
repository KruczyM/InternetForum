package database

import (
	"database/sql"
	"os"
	"strings"
)

func RunMigrations(db *sql.DB) error {
	files := []string{
		"internal/db/migrations/001_init.sql",
		"internal/db/migrations/002_add_comment_parent.sql",
	}

	for _, file := range files {
		query, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		if _, err := db.Exec(string(query)); err != nil {
			// Ignore migration error when column already exists (idempotent runs)
			if strings.Contains(err.Error(), "duplicate column name") || strings.Contains(err.Error(), "already exists") {
				continue
			}
			return err
		}
	}

	return nil
}
