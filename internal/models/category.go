package models

import (
	"database/sql"
	"strings"
)

type CategoryOverview struct {
	ID           int
	Name         string
	Kind         string
	Description  string
	AvatarPath   string
	Slug         string
	PostCount    int
	CommentCount int
}

type CategoryModel struct {
	DB *sql.DB
}

func (m *CategoryModel) GetCategoryOverview() ([]CategoryOverview, error) {
	stmt := `
		SELECT
			c.id,
			c.name,
			c.kind,
			COALESCE(c.description, '') AS description,
			COALESCE(c.avatar_path, '') AS avatar_path,
			c.name as slug,
			COUNT(DISTINCT p.id) AS post_count,
			COUNT(cm.id) AS comment_count

		FROM categories c
		LEFT JOIN posts p
			ON p.post_type = c.name
		LEFT JOIN comments cm
			ON cm.post_id = p.id

		GROUP BY c.id
		ORDER BY c.name;
		`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []CategoryOverview

	for rows.Next() {
		var c CategoryOverview
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Kind,
			&c.Description,
			&c.AvatarPath,
			&c.Slug,
			&c.PostCount,
			&c.CommentCount,
		)
		if err != nil {
			return nil, err
		}
		c.Name = strings.ReplaceAll(c.Name, "_", " ")
		categories = append(categories, c)
	}

	return categories, rows.Err()
}
