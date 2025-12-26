package models

import (
	"database/sql"
	"fmt"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) GetAllPosts() ([]PostView, error) {
	stmt := `
    SELECT p.id, p.user_id, p.title, p.content, p.post_type, p.book_id, p.chapter, p.created_at, u.username,
    COALESCE(SUM(l.value), 0)
    FROM posts p
    LEFT JOIN users u ON p.user_id = u.id
    LEFT JOIN likes l ON p.id = l.target_id AND l.target_type = 'post'
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

		var bookID sql.NullInt64
		var chapter sql.NullString


		err := rows.Scan(
			&pv.Post.ID,
			&pv.Post.UserID,
			&pv.Post.Title,
			&pv.Post.Content,
            &pv.Post.PostType,
			&bookID,
			&chapter,
			&pv.Post.CreatedAt,
			&pv.AuthorName,
			&pv.LikeCount,
		)
		if err != nil {
			fmt.Println("Scan Error:", err)
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

func (m *PostModel) GetPost(id int) (*PostView, error) {
	stmt := `
    SELECT 
        p.id, 
        p.user_id, 
        p.title, 
        p.content, 
        p.post_type, 
        COALESCE(p.book_id, 0),
        COALESCE(p.chapter, ''),
        p.created_at, 
        u.username,
        (SELECT COUNT(*) FROM likes WHERE target_id = p.id AND target_type = 'post') as like_count
    FROM posts p
    LEFT JOIN users u ON p.user_id = u.id
    WHERE p.id = ?`

	row := m.DB.QueryRow(stmt, id)

	pv := &PostView{}

	err := row.Scan(
		&pv.Post.ID,
		&pv.Post.UserID,
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

	commentStmt := `
    SELECT 
        c.id, 
        c.content, 
        c.created_at, 
        c.user_id, 
        u.username,
        (SELECT COUNT(*) FROM likes WHERE target_id = c.id AND target_type = 'comment')
    FROM comments c
    LEFT JOIN users u ON c.user_id = u.id
    WHERE c.post_id = ?
    ORDER BY c.created_at ASC`

	rows, err := m.DB.Query(commentStmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c Comment

		err = rows.Scan(&c.ID, &c.Content, &c.CreatedAt, &c.UserID, &c.UserName, &c.LikeCount)
		if err != nil {
			return nil, err
		}
        pv.Comments = append(pv.Comments, c)
    }
    return pv, nil
}

func (m *PostModel) InsertPost(userID string, title, content, postType string, bookID *int, chapter *string) (int, error) {
	stmt := `
	INSERT INTO posts (user_id, title, content, post_type, book_id, chapter, created_at)
	VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`

	result, err := m.DB.Exec(stmt, userID, title, content, postType, bookID, chapter)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *PostModel) InsertComment(postID int, userID string, content string) error {
	stmt := `INSERT INTO comments (post_id, user_id, content, created_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP)`

	_, err := m.DB.Exec(stmt, postID, userID, content)
	if err != nil {
		return err
	}
	return nil
}

func (m *PostModel) DeletePost(id int) error {
	stmt := `DELETE FROM posts WHERE id = ?`

	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *PostModel) UpdatePost(id int, title, content string) error {
	stmt := `UPDATE posts SET title = ?, content = ? WHERE id = ?`

	_, err := m.DB.Exec(stmt, title, content, id)
	if err != nil {
		return err
	}
	return nil
}

func (m* PostModel) ToggleLike(userID string, postID int) error {
	stmt := `SELECT value FROM likes WHERE user_id = ? AND target_type = 'post' AND target_id = ?`

	var value int
	err := m.DB.QueryRow(stmt, userID, postID).Scan(&value)

	if err == sql.ErrNoRows {
		stmt = `INSERT INTO likes (user_id, target_type, target_id, value) VALUES (?, 'post', ?, 1)`
		_, err = m.DB.Exec(stmt, userID, postID)
		return err
	} else if err != nil {
		return err
	}

	if value == 1 {
		stmt := `DELETE FROM likes WHERE user_id = ? AND target_type = 'post' AND target_id = ?`
		_, err = m.DB.Exec(stmt, userID, postID)

	} else {
		stmt := `UPDATE likes SET value = 1 WHERE user_id = ? AND target_type = 'post' AND target_id = ?`
		_, err = m.DB.Exec(stmt, userID, postID)
	}
	return err
}

func (m *PostModel) GetLikeCount(postID int) (int, error) {
	stmt := `SELECT COUNT(*) FROM likes WHERE target_type = 'post' AND target_id = ? AND value = 1`
	var count int
	err := m.DB.QueryRow(stmt, postID).Scan(&count)
	return count, err
}

func (m *PostModel) DeleteComment(id int, userID string) error {
	stmt := `DELETE FROM comments WHERE id = ? AND user_id = ?`

	result, err := m.DB.Exec(stmt, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (m *PostModel) LikeComment(commentID int, userID string) error {
	stmt := `SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = ? AND target_id = ? AND target_type = 'comment')`
	var exists bool
	err := m.DB.QueryRow(stmt, userID, commentID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = m.DB.Exec(`DELETE FROM likes WHERE user_id = ? AND target_id = ? AND target_type = 'comment'`, userID, commentID)
	} else {
		_, err = m.DB.Exec(`INSERT INTO likes (user_id, target_id, target_type, value) VALUES (?, ?, 'comment', 1)`, userID, commentID)
	}
	return err
}