package service

import (
       "context"
       "log"

       "github.com/redis/go-redis/v9"
       "github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/cache"
       "github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/db"
)

func GetTimelinePosts(ctx context.Context, userID int64) ([]db.Post, error) {
       posts, err := db.GetTimelinePosts(ctx, userID)
       if err != nil {
               log.Printf("Timeline DB error: %v", err)
       }
       return posts, err
}

func GetTimelineWithCache(ctx context.Context, userID int64) ([]db.Post, error) {
       posts, err := cache.GetTimelineFromCache(userID)
       if err != nil && err != redis.Nil {
               log.Printf("Redis GetTimelineFromCache error: %v", err)
               return nil, err
       }
	if posts != nil {
		return posts, nil
	}

	// Cache miss — Fetch from DB
       posts, err = db.GetTimelinePosts(ctx, userID)
       if err != nil {
               log.Printf("DB GetTimelinePosts error: %v", err)
               return nil, err
       }

	_ = cache.SetTimelineToCache(userID, posts)

	return posts, nil
}
