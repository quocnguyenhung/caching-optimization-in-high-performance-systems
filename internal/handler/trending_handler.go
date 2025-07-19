package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/service"
)

func GetTrending(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := int64(10) // Default top 10
	if limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 64); err == nil {
			limit = l
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	ids, hit, err := service.GetTrendingPosts(ctx, limit)
	if err != nil {
		http.Error(w, "Failed to fetch trending posts", http.StatusInternalServerError)
		return
	}

	if hit {
		w.Header().Set("X-Cache", "HIT")
	} else {
		w.Header().Set("X-Cache", "MISS")
	}

	json.NewEncoder(w).Encode(ids)
}
