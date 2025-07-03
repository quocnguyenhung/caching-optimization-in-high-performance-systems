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

func GetTimelineFromDB(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID).(int64)
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	posts, err := service.GetTimelinePosts(ctx, userID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "Timeline request timed out", http.StatusGatewayTimeout)
		} else {
			http.Error(w, "Failed to fetch timeline", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(posts)
}

func GetTimelineWithCache(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID).(int64)
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	posts, err := service.GetTimelineWithCache(ctx, userID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "Timeline request timed out", http.StatusGatewayTimeout)
		} else {
			http.Error(w, "Failed to fetch timeline with cache", http.StatusInternalServerError)
		}
		return
	}

	if posts == nil {
		posts = []db.Post{}
	}

	w.WriteHeader(http.StatusOK) // ✅ Make sure status code is 200
	json.NewEncoder(w).Encode(posts)
}
