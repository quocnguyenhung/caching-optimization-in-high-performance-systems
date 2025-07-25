package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/config"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/db"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/metrics"
)

type TTLStrategy string

const (
	ConventionalTTL  TTLStrategy = "conventional"
	InvertedTTL      TTLStrategy = "inverted"
	trendingCacheKey             = "trending"
)

func getTrendingCacheKey(limit int64) string {
	return fmt.Sprintf("trending:top:%d", limit)
}

func CacheTrendingPostsWithStrategy(limit int64, strategy TTLStrategy) ([]int64, error) {
	ctx := context.Background()
	cacheKey := getTrendingCacheKey(limit)

	val, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		fmt.Printf("[CACHE HIT] Key: %s, Strategy: %s\n", cacheKey, strategy)
		var ids []int64
		if jsonErr := json.Unmarshal([]byte(val), &ids); jsonErr == nil {
			metrics.IncHit(string(strategy))
			return ids, nil
		}
	}

	fmt.Printf("[CACHE MISS] Key: %s, Strategy: %s\n", cacheKey, strategy)
	metrics.IncMiss(string(strategy))

	ids, err := db.GetTopTrendingFromDB(limit)
	if err != nil {
		return nil, err
	}

	ttl := determineTTLFromStrategy(ids, strategy)
	fmt.Printf("[CACHE SET] Key: %s, TTL: %s, Strategy: %s\n", cacheKey, ttl, strategy)
	data, _ := json.Marshal(ids)
	_ = config.RedisClient.Set(ctx, cacheKey, data, ttl).Err()

	return ids, nil
}

func determineTTLFromStrategy(postIDs []int64, strategy TTLStrategy) time.Duration {
	if len(postIDs) == 0 {
		return time.Minute
	}

	ctx := context.Background()
	topID := postIDs[0]
	score, err := config.RedisClient.ZScore(ctx, trendingKey, fmt.Sprintf("%d", topID)).Result()
	if err != nil {
		score = 0
	}

	switch strategy {
	case ConventionalTTL:
		switch {
		case score >= 100:
			return time.Hour
		case score >= 10:
			return 10 * time.Minute
		default:
			return time.Minute
		}
	case InvertedTTL:
		switch {
		case score >= 100:
			return time.Minute
		case score >= 10:
			return 5 * time.Minute
		default:
			return time.Hour
		}
	default:
		return 5 * time.Minute
	}
}
