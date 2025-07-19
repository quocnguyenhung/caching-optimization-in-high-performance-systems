package service

import (
	"context"
	"log"

	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/cache"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/config"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/db"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/pkg/utils"
	"github.com/redis/go-redis/v9"
)

func GetTimelinePosts(ctx context.Context, userID int64) ([]db.Post, error) {
	posts, err := db.GetTimelinePosts(ctx, userID)
	if err != nil {
		log.Printf("Timeline DB error: %v", err)
	}
	return posts, err
}

func GetTimeline(ctx context.Context, userID int64) ([]db.Post, bool, error) {
	if config.EnableCaching && utils.UseCache(userID) {
		posts, hit, err := cache.GetTimelineFromCache(userID)
		if err != nil && err != redis.Nil {
			log.Printf("Redis GetTimelineFromCache error: %v", err)
			return nil, false, err
		}
		if hit {
			return posts, true, nil
		}

		posts, err = db.GetTimelinePosts(ctx, userID)
		if err != nil {
			log.Printf("DB GetTimelinePosts error: %v", err)
			return nil, false, err
		}

		_ = cache.SetTimelineToCache(userID, posts)
		return posts, false, nil
	}

	posts, err := db.GetTimelinePosts(ctx, userID)
	return posts, false, err
}
