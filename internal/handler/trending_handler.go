package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/cache"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/db"
)

func GetTrendingWithCache(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := int64(10) // Default top 10
	if limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 64); err == nil {
			limit = l
		}
	}

	ids, err := cache.GetTopTrendingPosts(limit)
	if err != nil {
		http.Error(w, "Failed to fetch trending posts", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ids)
}

func GetTrendingFromDB(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := int64(10)
	if limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 64); err == nil {
			limit = l
		}
	}

	ids, err := db.GetTopTrendingFromDB(limit)
	if err != nil {
		http.Error(w, "Failed to fetch trending from DB", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ids)
}
func TrendingWithTTLStrategy(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")
	limitStr := r.URL.Query().Get("limit")

	var limit int64 = 10
	if l, err := strconv.ParseInt(limitStr, 10, 64); err == nil {
		limit = l
	}

	strategy := cache.TTLStrategy(mode)
	start := time.Now()

	ids, err := cache.CacheTrendingPostsWithStrategy(limit, strategy)
	if err != nil {
		http.Error(w, "Failed to fetch trending posts", http.StatusInternalServerError)
		return
	}
	duration := time.Since(start)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"strategy": strategy,
		"took":     duration.String(),
		"data":     ids,
	})
}
