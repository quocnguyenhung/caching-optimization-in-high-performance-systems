package handler

import (
	"context"
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

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	ids, err := db.GetTopTrendingFromDB(ctx, limit)
	if err != nil {
		http.Error(w, "Failed to fetch trending from DB", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ids)
}
