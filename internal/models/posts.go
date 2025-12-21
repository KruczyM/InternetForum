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

func (m *PostModel) Get(id int) (*PostView, error) {
	stmt := `SELECT p.id, p.title, p.content, p.post_type, p.book_id, p.chapter, p.created_at, u.username,
	COALESCE(SUM(v.value), 0)
	FROM posts p
	LEFT JOIN users u ON p.user_id = u.id
	LEFT JOIN votes v ON p.id = v.target.id AND v.target_type = 'post
	WHERE p.id = ?
	GROUP BY p.id, u.username`

	row := m.DB.QueryRow(stmt, id)

	pv := &PostView{}

	err := row.Scan(
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

	pv.FormattedDate = pv.Post.CreatedAt.Format("Jul 09, 1990 at 5:04 PM")

	commentStmt := `
	SELECT c.id, c.content, c.created_at, u.username
	FROM comments c
	LEFT JOIN users u ON c.user_id u.id
	WHERE c.post_id = ?
	ORDER BY c.created_at ASC`

	rows, err := m.DB.Query(commentStmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c Comment

		err = rows.Scan(&c.ID, &c.Content, &c.CreatedAt, &c.UserID)
		if err != nil {
			return nil, err
		}
        pv.Comments = append(pv.Comments, c)
    }
    return pv, nil
}

