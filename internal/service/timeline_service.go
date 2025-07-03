package service

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/cache"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/db"
)

func GetTimelinePosts(ctx context.Context, userID int64) ([]db.Post, error) {
	posts, err := db.GetTimelinePosts(ctx, userID)
	if err != nil {
		fmt.Println("Timeline DB error:", err) // <-- add this line
	}
	return posts, err
}

func GetTimelineWithCache(ctx context.Context, userID int64) ([]db.Post, error) {
	posts, err := cache.GetTimelineFromCache(userID)
	if err != nil && err != redis.Nil {
		fmt.Println("Redis GetTimelineFromCache error:", err)
		return nil, err
	}
	if posts != nil {
		return posts, nil
	}

	// Cache miss — Fetch from DB
	posts, err = db.GetTimelinePosts(ctx, userID)
	if err != nil {
		fmt.Println("DB GetTimelinePosts error:", err) // <-- Add this print
		return nil, err
	}

	_ = cache.SetTimelineToCache(userID, posts)

	return posts, nil
}
