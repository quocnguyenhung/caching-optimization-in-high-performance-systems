package metrics

import (
	"sync/atomic"
)

type TTLMetrics struct {
	Hits   uint64
	Misses uint64
}

var (
	ConventionalStats TTLMetrics
	InvertedStats     TTLMetrics
)

func IncHit(strategy string) {
	switch strategy {
	case "conventional":
		atomic.AddUint64(&ConventionalStats.Hits, 1)
	case "inverted":
		atomic.AddUint64(&InvertedStats.Hits, 1)
	}
}

func IncMiss(strategy string) {
	switch strategy {
	case "conventional":
		atomic.AddUint64(&ConventionalStats.Misses, 1)
	case "inverted":
		atomic.AddUint64(&InvertedStats.Misses, 1)
	}
}

func GetStats(strategy string) (hits, misses uint64, hitRate float64) {
	switch strategy {
	case "conventional":
		hits = atomic.LoadUint64(&ConventionalStats.Hits)
		misses = atomic.LoadUint64(&ConventionalStats.Misses)
	case "inverted":
		hits = atomic.LoadUint64(&InvertedStats.Hits)
		misses = atomic.LoadUint64(&InvertedStats.Misses)
	}
	total := hits + misses
	if total == 0 {
		return hits, misses, 0.0
	}
	hitRate = float64(hits) / float64(total)
	return
}
