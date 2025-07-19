package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/config"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/db"
	"github.com/redis/go-redis/v9"
)

const timelineTTL = 5 * time.Minute // Default cache TTL

func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

func maxDuration(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}

// getAdaptiveTTL calculates TTL based on post volume when adaptive mode is enabled.
func getAdaptiveTTL(postCount int) time.Duration {
	if !config.EnableAdaptiveTTL {
		return timelineTTL
	}

	rate := postCount
	var newTTL time.Duration

	switch {
	case rate > config.TimelineHighThreshold:
		newTTL = minDuration(config.TimelineMaxTTL, 2*config.TimelineMinTTL)
	case rate < config.TimelineLowThreshold:
		newTTL = maxDuration(config.TimelineMinTTL, config.TimelineMinTTL/2)
	default:
		return timelineTTL
	}

	return newTTL
}

func GetTimelineFromCache(userID int64) ([]db.Post, bool, error) {
	key := fmt.Sprintf("timeline:%d", userID)

	ctx := context.Background()
	val, err := config.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		// cache miss
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	var posts []db.Post
	err = json.Unmarshal([]byte(val), &posts)
	if err != nil {
		return nil, false, err
	}

	return posts, true, nil
}

func SetTimelineToCache(userID int64, posts []db.Post) error {
	key := fmt.Sprintf("timeline:%d", userID)

	data, err := json.Marshal(posts)
	if err != nil {
		return err
	}

	ctx := context.Background()
	ttl := getAdaptiveTTL(len(posts))
	return config.RedisClient.Set(ctx, key, data, ttl).Err()
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
	ttl := getAdaptiveTTL(len(posts))
	return config.RedisClient.Set(ctx, key, data, ttl).Err()
}
