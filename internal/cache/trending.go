package cache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/config"
)

const trendingKey = "trending"

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
