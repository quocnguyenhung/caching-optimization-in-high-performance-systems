package handler

import (
	"encoding/json"
	"net/http"

	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/db"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/middleware"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/service"
)

func GetTimelineFromDB(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID).(int64)

	posts, err := service.GetTimelinePosts(userID)
	if err != nil {
		http.Error(w, "Failed to fetch timeline", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(posts)
}

func GetTimelineWithCache(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID).(int64)

	posts, err := service.GetTimelineWithCache(userID)
	if err != nil {
		http.Error(w, "Failed to fetch timeline with cache", http.StatusInternalServerError)
		return
	}

	if posts == nil {
		posts = []db.Post{}
	}

	w.WriteHeader(http.StatusOK) // ✅ Make sure status code is 200
	json.NewEncoder(w).Encode(posts)
}
