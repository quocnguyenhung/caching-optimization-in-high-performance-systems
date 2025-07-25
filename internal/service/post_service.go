package service

import (
	"errors"
	"time"

	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/cache"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/config"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/db"
)

func CreatePost(userID int64, content string) error {
	postID, err := db.CreatePost(userID, content)
	if err != nil {
		return err
	}

	post := db.Post{
		ID:        postID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if config.EnableCaching {
		// Write-through caching
		_ = cache.SetPostToCache(post)

		// Fan-out to followers
		followerIDs, err := db.GetFollowers(userID)
		if err == nil {
			for _, followerID := range followerIDs {
				_ = cache.PushPostToTimelineCache(followerID, post)
			}
		}
	}

	return nil
}

func FollowUser(followerID, followedID int64) error {
	if followerID == followedID {
		return errors.New("you cannot follow yourself")
	}

	exists, err := db.CheckFollowExists(followerID, followedID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("already following this user")
	}

	return db.FollowUser(followerID, followedID)
}

func LikePost(userID, postID int64) error {
	err := db.LikePost(userID, postID)
	if err != nil {
		return err
	}

	// Best-effort: Update trending cache asynchronously
	go func() {
		_ = cache.IncrementTrendingScore(postID)
	}()

	return nil
}
