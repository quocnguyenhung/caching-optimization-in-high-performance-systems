package config

import (
	"os"
	"time"
)

// EnableCaching controls whether Redis should be used.
var EnableCaching = os.Getenv("CACHE_ENABLED") != "false"

// EnableAdaptiveTTL toggles dynamic TTL calculation for timelines.
var EnableAdaptiveTTL = os.Getenv("ADAPTIVE_TTL_ENABLED") == "true"

// Timeline cache TTL settings used when adaptive mode is enabled.
const (
	TimelineMinTTL        = time.Minute
	TimelineMaxTTL        = 10 * time.Minute
	TimelineLowThreshold  = 10
	TimelineHighThreshold = 40
)
