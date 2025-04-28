package db

import "time"

type User struct {
	ID        int64     `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Password  string    `db:"password,omitempty"` // don't expose password in JSON
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Post struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"user_id"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Follow struct {
	ID         int64     `db:"id" json:"id"`
	FollowerID int64     `db:"follower_id" json:"follower_id"`
	FollowedID int64     `db:"followed_id" json:"followed_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type UserProfile struct {
	ID             int64  `json:"id"`
	Username       string `json:"username"`
	FollowerCount  int64  `json:"follower_count"`
	FollowingCount int64  `json:"following_count"`
	PostCount      int64  `json:"post_count"`
}
