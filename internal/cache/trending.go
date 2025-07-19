package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/config"
	"github.com/redis/go-redis/v9"
)

const trendingKey = "trending"
const trendingListTTL = 5 * time.Minute

func IncrementTrendingScore(postID int64) error {
	ctx := context.Background()
	return config.RedisClient.ZIncrBy(ctx, trendingKey, 1, fmt.Sprintf("%d", postID)).Err()
}

func GetTopTrendingPosts(limit int64) ([]int64, error) {
	ctx := context.Background()

	idsStr, err := config.RedisClient.ZRevRange(ctx, trendingKey, 0, limit-1).Result()
	if err != nil {
		return nil, err
	}

	var ids []int64
	for _, idStr := range idsStr {
		id, _ := strconv.ParseInt(idStr, 10, 64)
		ids = append(ids, id)
	}

	return ids, nil
}

// GetTopTrendingCached returns trending post IDs and whether it was a cache hit.
func GetTopTrendingCached(limit int64) ([]int64, bool, error) {
	cacheKey := fmt.Sprintf("trending:top:%d", limit)
	ctx := context.Background()

	val, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var ids []int64
		if err := json.Unmarshal([]byte(val), &ids); err == nil {
			return ids, true, nil
		}
	}
	if err != nil && err != redis.Nil {
		return nil, false, err
	}

	ids, err := GetTopTrendingPosts(limit)
	if err != nil {
		return nil, false, err
	}

	data, _ := json.Marshal(ids)
	_ = config.RedisClient.Set(ctx, cacheKey, data, trendingListTTL).Err()
	return ids, false, nil
}
