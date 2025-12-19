package models

import "database/sql"


// sets or updates a user's like/dislike for a post
func SetLike(db *sql.DB, userID, postID, value int) error {
	query := `INSERT INTO likes (user_id, target_type, target_id, value) VALUES (?, ?, ?, ?)
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