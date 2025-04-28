package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/config"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/db"
)

const postTTL = 24 * time.Hour

func SetPostToCache(post db.Post) error {
	key := fmt.Sprintf("post:%d", post.ID)

	data, err := json.Marshal(post)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return config.RedisClient.Set(ctx, key, data, postTTL).Err()
}

func GetPostFromCache(postID int64) (*db.Post, error) {
	key := fmt.Sprintf("post:%d", postID)

	val, err := config.RedisClient.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	var post db.Post
	if err := json.Unmarshal([]byte(val), &post); err != nil {
		return nil, err
	}

	return &post, nil
}
