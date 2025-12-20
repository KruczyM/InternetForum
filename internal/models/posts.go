package models

import (
	"database/sql"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) GetAll() ([]PostView, error) {
	stmt := `SELECT
			p.ID, p.title, p.content, p.post_type, p.book_id, p.chapter, p.created_at,
			u.username,
			COALESCE(SUM(v.value), 0)
			FROM posts p
			LEFT JOIN users u ON p.user:id = u.id
			LEFT JOIN votes v ON p.id = v.target_id AND v.target_type = 'post'
			GROUP BY p.id, u.username
			ORDER BY p.created_at DESC`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []PostView

	for rows.Next() {
		var pv PostView

		err := rows.Scan(
			&pv.Post.ID,
			&pv.Post.Title,
			&pv.Post.Content,
            &pv.Post.PostType,
			&pv.Post.BookID,
			&pv.Post.Chapter,
			&pv.Post.CreatedAt,
			&pv.AuthorName,
			&pv.LikeCount,
		)
		if err != nil {
			return nil, err
		}

		pv.FormattedDate = pv.Post.CreatedAt.Format("Jan 02, 2006 at 3:04 PM")

		posts = append(posts, pv)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}