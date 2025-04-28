package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/service"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	userIDStr := pathParts[2]
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	profile, err := service.GetUserProfile(userID)
	if err != nil {
		http.Error(w, "Failed to get profile", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(profile)
}
