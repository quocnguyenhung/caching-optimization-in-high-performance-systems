package service

import (
	"context"

	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/cache"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/db"
)

func GetUserProfile(ctx context.Context, userID int64) (*db.UserProfile, error) {
	// First try cache
	profile, err := cache.GetProfileCache(userID)
	if err == nil && profile != nil {
		return profile, nil
	}

	// Fallback to DB
	profile, err = db.GetUserProfileFromDB(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update cache
	_ = cache.SetProfileCache(*profile)

	return profile, nil
}
