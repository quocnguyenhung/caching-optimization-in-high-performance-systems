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

const timelineTTL = 5 * time.Minute // Cache expires after 5 minutes

func GetTimelineFromCache(userID int64) ([]db.Post, error) {
	key := fmt.Sprintf("timeline:%d", userID)
	ctx := context.Background()

	vals, err := config.RedisClient.LRange(ctx, key, 0, 49).Result()
	if err == redis.Nil {
		return nil, nil // cache miss
	}
	if err != nil {
		return nil, err
	}

	var posts []db.Post
	for _, val := range vals {
		var post db.Post
		if err := json.Unmarshal([]byte(val), &post); err == nil {
			posts = append(posts, post)
		}
	}

	return posts, nil
}

func SetTimelineToCache(userID int64, posts []db.Post) error {
	key := fmt.Sprintf("timeline:%d", userID)
	ctx := context.Background()

	pipe := config.RedisClient.TxPipeline()

	for i := len(posts) - 1; i >= 0; i-- {
		data, err := json.Marshal(posts[i])
		if err != nil {
			continue
		}
		pipe.LPush(ctx, key, data)
	}

	pipe.LTrim(ctx, key, 0, 49)
	pipe.Expire(ctx, key, timelineTTL)

	_, err := pipe.Exec(ctx)
	return err
}

func PushPostToTimelineCache(userID int64, post db.Post) error {
	key := fmt.Sprintf("timeline:%d", userID)
	ctx := context.Background()

	data, err := json.Marshal(post)
	if err != nil {
		return err
	}

	pipe := config.RedisClient.Pipeline()
	pipe.LPush(ctx, key, data)
	pipe.LTrim(ctx, key, 0, 49)
	pipe.Expire(ctx, key, timelineTTL)

	_, err = pipe.Exec(ctx)
	return err
}
