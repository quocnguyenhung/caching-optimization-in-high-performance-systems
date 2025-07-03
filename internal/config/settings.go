package config

import (
	"os"
	"time"
)

var EnableCaching = true

// EnableAdaptiveTTL toggles dynamic TTL calculation for timelines.
// Controlled via the ENABLE_ADAPTIVE_TTL environment variable.
var EnableAdaptiveTTL = os.Getenv("ENABLE_ADAPTIVE_TTL") == "true"

// Timeline cache TTL settings used when adaptive mode is enabled.
const (
	TimelineMinTTL        = time.Minute
	TimelineMaxTTL        = 10 * time.Minute
	TimelineLowThreshold  = 10
	TimelineHighThreshold = 40
)
