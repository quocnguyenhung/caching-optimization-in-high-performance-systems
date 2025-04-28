package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/config"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/db"
)

const timelineTTL = 5 * time.Minute // Cache expires after 5 mins

func GetTimelineFromCache(userID int64) ([]db.Post, error) {
	key := fmt.Sprintf("timeline:%d", userID)

	ctx := context.Background()
	val, err := config.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		// cache miss
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var posts []db.Post
	err = json.Unmarshal([]byte(val), &posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func SetTimelineToCache(userID int64, posts []db.Post) error {
	key := fmt.Sprintf("timeline:%d", userID)

	data, err := json.Marshal(posts)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return config.RedisClient.Set(ctx, key, data, timelineTTL).Err()
}

func PushPostToTimelineCache(userID int64, post db.Post) error {
	key := fmt.Sprintf("timeline:%d", userID)

	ctx := context.Background()

	// Get existing timeline
	val, err := config.RedisClient.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	var posts []db.Post
	if val != "" {
		_ = json.Unmarshal([]byte(val), &posts)
	}

	// Prepend the new post
	posts = append([]db.Post{post}, posts...)

	// Keep only latest 50 posts
	if len(posts) > 50 {
		posts = posts[:50]
	}

	// Save updated timeline back
	data, _ := json.Marshal(posts)
	return config.RedisClient.Set(ctx, key, data, timelineTTL).Err()
}
