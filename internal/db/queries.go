package db

import (
	"context"
	"time"

	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/config"
)

// timeout applied to all DB queries to avoid exhausting connections
const dbTimeout = 3 * time.Second

// CreateUser inserts a new user into the DB
func CreateUser(ctx context.Context, username, hashedPassword string) error {
	c, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	_, err := config.DB.ExecContext(c, "INSERT INTO users (username, password) VALUES ($1, $2)", username, hashedPassword)
	return err
}

// GetUserByUsername fetches a user by username
func GetUserByUsername(ctx context.Context, username string) (*User, error) {
	c, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	var user User
	err := config.DB.QueryRowContext(c, "SELECT id, username, password, created_at FROM users WHERE username=$1", username).
		Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreatePost inserts a new post into the DB
func CreatePost(ctx context.Context, userID int64, content string) (int64, error) {
	query := `
               INSERT INTO posts (user_id, content, created_at)
               VALUES ($1, $2, now())
               RETURNING id;
       `

	c, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	var postID int64
	err := config.DB.QueryRowContext(c, query, userID, content).Scan(&postID)
	if err != nil {
		return 0, err
	}

	return postID, nil
}

// GetTimelinePosts fetches posts for a user and people they follow
func GetTimelinePosts(ctx context.Context, userID int64) ([]Post, error) {
	query := `
	SELECT posts.id, posts.user_id, posts.content, posts.created_at
	FROM posts
	WHERE posts.user_id = $1
	OR posts.user_id IN (
		SELECT followed_id FROM follows WHERE follower_id = $1
	)
	ORDER BY posts.created_at DESC
	LIMIT 50;`

	c, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	rows, err := config.DB.QueryContext(c, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.Content, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

// FollowUser adds a follow relationship
func FollowUser(ctx context.Context, followerID, followedID int64) error {
	c, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	_, err := config.DB.ExecContext(c, "INSERT INTO follows (follower_id, followed_id) VALUES ($1, $2)", followerID, followedID)
	return err
}

// CheckFollowExists checks if a follow relationship already exists
func CheckFollowExists(ctx context.Context, followerID, followedID int64) (bool, error) {
	c, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	var exists bool
	err := config.DB.QueryRowContext(
		c,
		"SELECT EXISTS(SELECT 1 FROM follows WHERE follower_id=$1 AND followed_id=$2)",
		followerID, followedID,
	).Scan(&exists)
	return exists, err
}

func GetFollowers(ctx context.Context, userID int64) ([]int64, error) {
	query := `
       SELECT follower_id FROM follows WHERE followed_id = $1;
       `

	c, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	rows, err := config.DB.QueryContext(c, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followerIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		followerIDs = append(followerIDs, id)
	}

	return followerIDs, nil
}

func LikePost(ctx context.Context, userID, postID int64) error {
	query := `
               INSERT INTO likes (user_id, post_id, created_at)
               VALUES ($1, $2, now())
               ON CONFLICT DO NOTHING;
       `
	c, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	_, err := config.DB.ExecContext(c, query, userID, postID)
	if err != nil {
		return err
	}

	_, err = config.DB.ExecContext(c, "UPDATE posts SET likes = likes + 1 WHERE id=$1", postID)
	return err
}

func GetTopTrendingFromDB(ctx context.Context, limit int64) ([]int64, error) {
	query := `
	SELECT post_id
	FROM (
	    SELECT post_id, COUNT(*) as like_count
	    FROM likes
	    GROUP BY post_id
	) as likes_counts
	ORDER BY like_count DESC
	LIMIT $1;
	`

	c, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	rows, err := config.DB.QueryContext(c, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var postIDs []int64
	for rows.Next() {
		var postID int64
		if err := rows.Scan(&postID); err != nil {
			return nil, err
		}
		postIDs = append(postIDs, postID)
	}

	return postIDs, nil
}

func GetUserProfileFromDB(ctx context.Context, userID int64) (*UserProfile, error) {
	query := `
	SELECT id, username,
		(SELECT COUNT(*) FROM follows WHERE followed_id = users.id) as follower_count,
		(SELECT COUNT(*) FROM follows WHERE follower_id = users.id) as following_count,
		(SELECT COUNT(*) FROM posts WHERE user_id = users.id) as post_count
	FROM users
	WHERE id = $1;
	`

	var profile UserProfile
	c, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	err := config.DB.QueryRowContext(c, query, userID).Scan(
		&profile.ID,
		&profile.Username,
		&profile.FollowerCount,
		&profile.FollowingCount,
		&profile.PostCount,
	)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}
