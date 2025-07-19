package service

import (
	"context"

	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/cache"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/config"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/db"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/pkg/utils"
	"github.com/redis/go-redis/v9"
)

func GetUserProfile(ctx context.Context, userID int64) (*db.UserProfile, bool, error) {
	if config.EnableCaching && utils.UseCache(userID) {
		profile, hit, err := cache.GetProfileCache(userID)
		if err == nil && hit {
			return profile, true, nil
		}
		if err != nil && err != redis.Nil {
			return nil, false, err
		}
	}

	// Fallback to DB
	profile, err := db.GetUserProfileFromDB(ctx, userID)
	if err != nil {
		return nil, false, err
	}

	// Update cache
	if config.EnableCaching {
		_ = cache.SetProfileCache(*profile)
	}

	return profile, false, nil
}
