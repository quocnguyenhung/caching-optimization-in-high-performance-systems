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

const profileTTL = 5 * time.Minute

func SetProfileCache(profile db.UserProfile) error {
	key := fmt.Sprintf("profile:%d", profile.ID)

	data, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return config.RedisClient.Set(ctx, key, data, profileTTL).Err()
}

func GetProfileCache(userID int64) (*db.UserProfile, bool, error) {
	key := fmt.Sprintf("profile:%d", userID)

	val, err := config.RedisClient.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	var profile db.UserProfile
	if err := json.Unmarshal([]byte(val), &profile); err != nil {
		return nil, false, err
	}

	return &profile, true, nil
}
