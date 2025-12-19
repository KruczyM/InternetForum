package models

import "database/sql"


// sets or updates a user's like/dislike for a post
func SetLike(db *sql.DB, userID, postID, value int) error {
	query := `
	INSERT INTO likes (user_id, target_type, target_id, value) 
	VALUES (?, ?, ?, ?)
		ON CONFLICT(user_id, target_type, target_id) DO UPDATE SET value=excluded.value`
	_, err := db.Exec(query, userID, postID, value)
	return err
}

// returns the number of likes and dislikes for a post
func GetPostLikes(db *sql.DB, postID int) (likes int, dislikes int, err error) {

	// COALESCE is used to return 0 instead of NULL	
	// sum all when value = 1, then +1 else +0 
	query := `
		SELECT
			COALESCE(SUM(CASE WHEN value = 1 THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN value = -1 THEN 1 ELSE 0 END), 0)
		FROM likes
		WHERE target_type = 'post' AND target_id = ?
	`
	err = db.QueryRow(query, postID).Scan(&likes, &dislikes)
	return
}

func CreatePost(db *sql.DB, userID string, title string, content string) error {
	query := `
	INSERT INTO posts (user_id, title, content) 
	VALUES (?, ?, ?)
	`
	_, err := db.Exec(query, userID, title, content)
	return err
}

func GetAllPosts(db *sql.DB) ([]Post, error) {
	query := `
		SELECT id, user_id, title, content, post_type, book_id, chapter, created_at
		FROM posts
		ORDER BY created_at DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var p Post
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.PostType,
			&p.BookID,
			&p.Chapter,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}