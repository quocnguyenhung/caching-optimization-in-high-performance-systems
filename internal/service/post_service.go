package service

import (
	"context"
	"errors"
	"time"

	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/cache"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/config"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/db"
)

func CreatePost(ctx context.Context, userID int64, content string) error {
	postID, err := db.CreatePost(ctx, userID, content)
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
		followerIDs, err := db.GetFollowers(ctx, userID)
		if err == nil {
			for _, followerID := range followerIDs {
				_ = cache.PushPostToTimelineCache(followerID, post)
			}
		}
	}

	return nil
}

func FollowUser(ctx context.Context, followerID, followedID int64) error {
	if followerID == followedID {
		return errors.New("you cannot follow yourself")
	}

	exists, err := db.CheckFollowExists(ctx, followerID, followedID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("already following this user")
	}

	return db.FollowUser(ctx, followerID, followedID)
}

func LikePost(ctx context.Context, userID, postID int64) error {
	err := db.LikePost(ctx, userID, postID)
	if err != nil {
		return err
	}

	// Best-effort: Update trending cache asynchronously
	go func() {
		_ = cache.IncrementTrendingScore(postID)
	}()

	return nil
}
