package service

import (
	"context"

	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/cache"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/config"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/db"
)

// GetTrendingPosts returns IDs of trending posts and whether they were served from cache.
func GetTrendingPosts(ctx context.Context, limit int64) ([]int64, bool, error) {
	if config.EnableCaching {
		ids, hit, err := cache.GetTopTrendingCached(limit)
		if err == nil {
			return ids, hit, nil
		}
		// fallback to DB in case of error
	}

	ids, err := db.GetTopTrendingFromDB(ctx, limit)
	return ids, false, err
}
