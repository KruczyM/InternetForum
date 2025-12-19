package database

import "database/sql"

// Like represents a user's like or dislike on a post
// Value: 1 for like, -1 for dislike
// Each (UserID, PostID) pair is unique

type Like struct {
	UserID int
	PostID int
	Value  int // 1 for like, -1 for dislike
}

// creates the likes table if it doesn't exist
func CreateLikesTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS likes (
		user_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		value INTEGER NOT NULL,
		PRIMARY KEY (user_id, post_id)
	);`
	_, err := db.Exec(query)
	return err
}

// sets or updates a user's like/dislike for a post
func SetLike(db *sql.DB, userID, postID, value int) error {
	query := `INSERT INTO likes (user_id, post_id, value) VALUES (?, ?, ?)
		ON CONFLICT(user_id, post_id) DO UPDATE SET value=excluded.value`
	_, err := db.Exec(query, userID, postID, value)
	return err
}

// returns the number of likes and dislikes for a post
func GetPostLikes(db *sql.DB, postID int) (likes int, dislikes int, err error) {
	query := `SELECT value, COUNT(*) FROM likes WHERE post_id = ? GROUP BY value`
	rows, err := db.Query(query, postID)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var value, count int
		if err := rows.Scan(&value, &count); err != nil {
			return 0, 0, err
		}
		if value == 1 {
			likes = count
		} else if value == -1 {
			dislikes = count
		}
	}
	return likes, dislikes, nil
}
