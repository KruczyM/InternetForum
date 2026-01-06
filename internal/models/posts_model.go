package models

import (
	"database/sql"
	"fmt"
	"strings"
)

type PostView struct {
	Post
	AuthorName    string
	BookTitle     string
	LikeCount     int
	DislikeCount  int
	CommentCount  int
	FormattedDate string
	Comments      []Comment
}

type PageData struct {
	Posts []PostView
}

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) GetAllPosts(category string, bookID int) ([]PostView, error) {
	stmt := `
	SELECT p.id, p.user_id, p.title, p.content, p.image_path, p.post_type, p.book_id,COALESCE(b.title, '') as book_title, p.chapter, p.created_at, u.username,
	COALESCE(SUM(CASE WHEN l.value = 1 THEN 1 ELSE 0 END), 0) as like_count,
	COALESCE(SUM(CASE WHEN l.value = -1 THEN 1 ELSE 0 END), 0) as dislike_count
	FROM posts p
	LEFT JOIN users u ON p.user_id = u.id
	LEFT JOIN books b ON p.book_id = b.id
	LEFT JOIN likes l ON p.id = l.target_id AND l.target_type = 'post'
   WHERE 1=1`

	var args []interface{}

	if category != "" {
		stmt += " AND p.post_type = ?"
		args = append(args, category)
	}

	if bookID > 0 {
		stmt += " AND p.book_id = ?"
		args = append(args, bookID)
	}

	stmt += ` GROUP BY p.id, username ORDER BY p.created_at DESC`

	rows, err := m.DB.Query(stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []PostView

	for rows.Next() {
		var pv PostView

		var bookIDNull sql.NullInt64
		var chapterNull sql.NullString

		err := rows.Scan(
			&pv.Post.ID,
			&pv.Post.UserID,
			&pv.Post.Title,
			&pv.Post.Content,
			&pv.Post.ImagePath,
			&pv.Post.PostType,
			&bookIDNull,
			&pv.BookTitle,
			&chapterNull,
			&pv.Post.CreatedAt,
			&pv.AuthorName,
			&pv.LikeCount,
			&pv.DislikeCount,
		)
		if err != nil {
			fmt.Println("Scan Error:", err)
			return nil, err
		}

		if bookIDNull.Valid {
			bID := int(bookIDNull.Int64)
			pv.Post.BookID = &bID
		}
		if chapterNull.Valid {
			chap := chapterNull.String
			pv.Post.Chapter = &chap
		}

		pv.FormattedDate = pv.Post.CreatedAt.Format("Jan 02, 2006 at 3:04 PM")

		posts = append(posts, pv)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (m *PostModel) GetPostsByUserID(userID string) ([]PostView, error) {
	stmt := `
	SELECT 
		p.id, p.user_id, p.title, p.content, p.image_path,
		p.post_type, p.book_id, p.chapter, p.created_at,
		u.username,
		COALESCE(SUM(l.value), 0)
	FROM posts p
	JOIN users u ON p.user_id = u.id
	LEFT JOIN likes l ON p.id = l.target_id AND l.target_type = 'post'
	WHERE p.user_id = ?
	GROUP BY p.id, u.username
	ORDER BY p.created_at DESC
	`

	rows, err := m.DB.Query(stmt, userID)
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
			&pv.Post.ImagePath,
			&pv.Post.PostType,
			&bookID,
			&chapter,
			&pv.Post.CreatedAt,
			&pv.AuthorName,
			&pv.LikeCount,
		)
		if err != nil {
			return nil, err
		}

		if bookID.Valid {
			id := int(bookID.Int64)
			pv.Post.BookID = &id
		}
		if chapter.Valid {
			c := chapter.String
			pv.Post.Chapter = &c
		}

		pv.FormattedDate = pv.Post.CreatedAt.Format("Jan 02, 2006 at 3:04 PM")
		posts = append(posts, pv)
	}

	return posts, rows.Err()
}

func (m *PostModel) GetPost(id int) (*PostView, error) {
	stmt := `
	SELECT 
		p.id, 
		p.user_id, 
		p.title, 
		p.content,
		p.image_path,
		p.post_type, 
		COALESCE(b.title, '') as book_title,
		p.book_id,
		COALESCE(p.chapter, '') as chapter,
		p.created_at, 
		u.username,
		(SELECT COUNT(*) FROM likes WHERE target_id = p.id AND target_type = 'post' AND value = 1) as like_count,
		(SELECT COUNT(*) FROM likes WHERE target_id = p.id AND target_type = 'post' AND value = -1) as dislike_count
	FROM posts p
	LEFT JOIN users u ON p.user_id = u.id
	LEFT JOIN books b ON p.book_id = b.id
	WHERE p.id = ?`

	row := m.DB.QueryRow(stmt, id)

	pv := &PostView{}

	err := row.Scan(
		&pv.Post.ID,
		&pv.Post.UserID,
		&pv.Post.Title,
		&pv.Post.Content,
		&pv.Post.ImagePath,
		&pv.Post.PostType,
		&pv.BookTitle,
		&pv.BookID,
		&pv.Post.Chapter,
		&pv.Post.CreatedAt,
		&pv.AuthorName,
		&pv.LikeCount,
		&pv.DislikeCount,
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
		(SELECT COUNT(*) FROM likes WHERE target_id = c.id AND target_type = 'comment' AND value = 1) as like_count,
		(SELECT COUNT(*) FROM likes WHERE target_id = c.id AND target_type = 'comment' AND value = -1) as dislike_count
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
		err = rows.Scan(&c.ID, &c.Content, &c.CreatedAt, &c.UserID, &c.UserName, &c.LikeCount, &c.DislikeCount)
		if err != nil {
			return nil, err
		}
		pv.Comments = append(pv.Comments, c)
	}
	return pv, nil
}

func (m *PostModel) InsertPost(userID string, title, content, imagePath string, postType string, bookID *int, chapter *string) (int, error) {
	stmt := `
	INSERT INTO posts (user_id, title, content, image_path, post_type, book_id, chapter, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`

	result, err := m.DB.Exec(stmt, userID, title, content, imagePath, postType, bookID, chapter)
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

// ToggleLike sets a like (1) for the post, or switches from dislike (-1) to like (1).
func (m *PostModel) ToggleLike(userID string, postID int) error {
	stmt := `SELECT value FROM likes WHERE user_id = ? AND target_type = 'post' AND target_id = ?`
	var value int
	err := m.DB.QueryRow(stmt, userID, postID).Scan(&value)
	if err == sql.ErrNoRows {
		// No reaction yet, insert like
		stmt = `INSERT INTO likes (user_id, target_type, target_id, value) VALUES (?, 'post', ?, 1)`
		_, err = m.DB.Exec(stmt, userID, postID)
		return err
	} else if err != nil {
		return err
	}
	if value == 1 {
		// Already liked, remove reaction
		stmt = `DELETE FROM likes WHERE user_id = ? AND target_type = 'post' AND target_id = ?`
		_, err = m.DB.Exec(stmt, userID, postID)
	} else {
		// Was disliked, switch to like
		stmt = `UPDATE likes SET value = 1 WHERE user_id = ? AND target_type = 'post' AND target_id = ?`
		_, err = m.DB.Exec(stmt, userID, postID)
	}
	return err
}

// ToggleDislike sets a dislike (-1) for the post, or switches from like (1) to dislike (-1).
func (m *PostModel) ToggleDislike(userID string, postID int) error {
	stmt := `SELECT value FROM likes WHERE user_id = ? AND target_type = 'post' AND target_id = ?`
	var value int
	err := m.DB.QueryRow(stmt, userID, postID).Scan(&value)
	if err == sql.ErrNoRows {
		// No reaction yet, insert dislike
		stmt = `INSERT INTO likes (user_id, target_type, target_id, value) VALUES (?, 'post', ?, -1)`
		_, err = m.DB.Exec(stmt, userID, postID)
		return err
	} else if err != nil {
		return err
	}
	if value == -1 {
		// Already disliked, remove reaction
		stmt = `DELETE FROM likes WHERE user_id = ? AND target_type = 'post' AND target_id = ?`
		_, err = m.DB.Exec(stmt, userID, postID)
	} else {
		// Was liked, switch to dislike
		stmt = `UPDATE likes SET value = -1 WHERE user_id = ? AND target_type = 'post' AND target_id = ?`
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

// ToggleLikeComment sets a like (1) for the comment, or switches from dislike (-1) to like (1).
func (m *PostModel) ToggleLikeComment(userID string, commentID int) error {
	stmt := `SELECT value FROM likes WHERE user_id = ? AND target_type = 'comment' AND target_id = ?`
	var value int
	err := m.DB.QueryRow(stmt, userID, commentID).Scan(&value)
	if err == sql.ErrNoRows {
		// No reaction yet, insert like
		stmt = `INSERT INTO likes (user_id, target_type, target_id, value) VALUES (?, 'comment', ?, 1)`
		_, err = m.DB.Exec(stmt, userID, commentID)
		return err
	} else if err != nil {
		return err
	}
	if value == 1 {
		// Already liked, remove reaction
		stmt = `DELETE FROM likes WHERE user_id = ? AND target_type = 'comment' AND target_id = ?`
		_, err = m.DB.Exec(stmt, userID, commentID)
	} else {
		// Was disliked, switch to like
		stmt = `UPDATE likes SET value = 1 WHERE user_id = ? AND target_type = 'comment' AND target_id = ?`
		_, err = m.DB.Exec(stmt, userID, commentID)
	}
	return err
}

// ToggleDislikeComment sets a dislike (-1) for the comment, or switches from like (1) to dislike (-1).
func (m *PostModel) ToggleDislikeComment(userID string, commentID int) error {
	stmt := `SELECT value FROM likes WHERE user_id = ? AND target_type = 'comment' AND target_id = ?`
	var value int
	err := m.DB.QueryRow(stmt, userID, commentID).Scan(&value)
	if err == sql.ErrNoRows {
		// No reaction yet, insert dislike
		stmt = `INSERT INTO likes (user_id, target_type, target_id, value) VALUES (?, 'comment', ?, -1)`
		_, err = m.DB.Exec(stmt, userID, commentID)
		return err
	} else if err != nil {
		return err
	}
	if value == -1 {
		// Already disliked, remove reaction
		stmt = `DELETE FROM likes WHERE user_id = ? AND target_type = 'comment' AND target_id = ?`
		_, err = m.DB.Exec(stmt, userID, commentID)
	} else {
		// Was liked, switch to dislike
		stmt = `UPDATE likes SET value = -1 WHERE user_id = ? AND target_type = 'comment' AND target_id = ?`
		_, err = m.DB.Exec(stmt, userID, commentID)
	}
	return err
}

func (m *PostModel) GetAllBooks() ([]Book, error) {
	stmt := `SELECT id, title FROM books ORDER BY title ASC`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book

	for rows.Next() {
		var b Book
		err = rows.Scan(&b.ID, &b.Title)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return books, nil

}

func (m *PostModel) SearchPosts(query, category string, bookID int, sort string) ([]PostView, error) {
	baseStmt := `
    SELECT 
        p.id, 
        p.user_id, 
        p.title, 
        p.content, 
        p.post_type, 
        p.book_id, 
        p.chapter, 
        p.created_at,
        u.username,
        COALESCE((
            SELECT SUM(l.value) 
            FROM likes l 
            WHERE l.target_id = p.id AND l.target_type='post'
        ), 0) AS likes,
        b.title AS book_title
    FROM posts p
    LEFT JOIN users u ON p.user_id = u.id
    LEFT JOIN books b ON p.book_id = b.id
    WHERE 1=1
`
	args := []interface{}{}

	// --- search query ---
	if strings.TrimSpace(query) != "" {
		baseStmt += " AND (p.title LIKE ? OR u.username LIKE ?)"
		likeQuery := "%" + query + "%"
		args = append(args, likeQuery, likeQuery)
	}

	// --- category filter ---
	if category != "" {
		baseStmt += " AND p.post_type = ?"
		args = append(args, category)
	}

	// --- book filter ---
	if bookID > 0 {
		baseStmt += " AND p.book_id = ?"
		args = append(args, bookID)
	}

	// --- sorting ---
	switch sort {
	case "popular":
		baseStmt += " ORDER BY likes DESC, p.created_at DESC"
	default:
		baseStmt += " ORDER BY p.created_at DESC"
	}

	rows, err := m.DB.Query(baseStmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []PostView
	for rows.Next() {
		var pv PostView
		var dbBookID sql.NullInt64
		var chapter sql.NullString
		var bookTitle sql.NullString

		if err := rows.Scan(
			&pv.Post.ID,
			&pv.Post.UserID,
			&pv.Post.Title,
			&pv.Post.Content,
			&pv.Post.PostType,
			&dbBookID,
			&chapter,
			&pv.Post.CreatedAt,
			&pv.AuthorName,
			&pv.LikeCount,
			&bookTitle,
		); err != nil {
			return nil, err
		}

		if dbBookID.Valid {
			id := int(dbBookID.Int64)
			pv.Post.BookID = &id
		}

		if chapter.Valid {
			pv.Post.Chapter = &chapter.String
		}

		if bookTitle.Valid {
			pv.BookTitle = bookTitle.String
		} else {
			pv.BookTitle = ""
		}

		pv.FormattedDate = pv.Post.CreatedAt.Format("Jan 02, 2006 at 3:04 PM")
		posts = append(posts, pv)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

type CommentView struct {
	ID        int
	Content   string
	PostID    int
	PostTitle string
}

func (m *PostModel) GetCommentsByUserID(userID string) ([]CommentView, error) {
	stmt := `
	SELECT 
		c.id, c.content,
		p.id, p.title
	FROM comments c
	JOIN posts p ON c.post_id = p.id
	WHERE c.user_id = ?
	ORDER BY c.created_at DESC
	`

	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []CommentView

	for rows.Next() {
		var cv CommentView

		err := rows.Scan(
			&cv.ID,
			&cv.Content,
			&cv.PostID,
			&cv.PostTitle,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, cv)
	}

	return comments, rows.Err()
}

type LikeView struct {
	TargetType string // "post" | "comment"
	TargetID   int
	PostID     int
	Title      string
	Content    string
	Value      int // 1 for like, -1 for dislike
}

func (m *PostModel) GetLikesByUserID(userID string) ([]LikeView, error) {
	stmt := `
	   -- liked POSTS
	   SELECT 
		   'post' AS target_type,
		   p.id AS target_id,
		   p.id AS post_id,
		   p.title,
		   p.content,
		   l.value
	   FROM likes l
	   JOIN posts p ON l.target_id = p.id
	   WHERE l.user_id = ? AND l.target_type = 'post'

	   UNION ALL

	   -- liked COMMENTS
	   SELECT
		   'comment' AS target_type,
		   c.id AS target_id,
		   p.id AS post_id,
		   p.title,
		   c.content,
		   l.value
	   FROM likes l
	   JOIN comments c ON l.target_id = c.id
	   JOIN posts p ON c.post_id = p.id
	   WHERE l.user_id = ? AND l.target_type = 'comment'

	   ORDER BY post_id DESC
	   `

	rows, err := m.DB.Query(stmt, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likes []LikeView

	for rows.Next() {
		var lv LikeView
		err := rows.Scan(
			&lv.TargetType,
			&lv.TargetID,
			&lv.PostID,
			&lv.Title,
			&lv.Content,
			&lv.Value,
		)
		if err != nil {
			return nil, err
		}
		likes = append(likes, lv)
	}

	return likes, rows.Err()
}
