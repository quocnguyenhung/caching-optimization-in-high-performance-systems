package handler

import (
	"encoding/json"
	"net/http"

	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/metrics"
)

func TTLStatsHandler(w http.ResponseWriter, r *http.Request) {
	strategy := r.URL.Query().Get("strategy")
	if strategy != "conventional" && strategy != "inverted" {
		http.Error(w, "Invalid strategy", http.StatusBadRequest)
		return
	}

	hits, misses, rate := metrics.GetStats(strategy)

	resp := map[string]interface{}{
		"strategy": strategy,
		"hits":     hits,
		"misses":   misses,
		"hit_rate": rate,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
