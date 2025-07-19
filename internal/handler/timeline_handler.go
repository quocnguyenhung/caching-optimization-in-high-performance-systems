package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/db"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/middleware"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/service"
)

func GetTimeline(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID).(int64)
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	posts, hit, err := service.GetTimeline(ctx, userID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "Timeline request timed out", http.StatusGatewayTimeout)
		} else {
			http.Error(w, "Failed to fetch timeline", http.StatusInternalServerError)
		}
		return
	}

	if posts == nil {
		posts = []db.Post{}
	}

	if hit {
		w.Header().Set("X-Cache", "HIT")
	} else {
		w.Header().Set("X-Cache", "MISS")
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}
