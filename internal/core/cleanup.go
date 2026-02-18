package core

import (
	"log/slog"
	"time"
)

const (
	cleanupInterval      = 10 * time.Minute
	rateLimiterMaxAge    = 1 * time.Hour
	posterCacheMaxAge    = 24 * time.Hour
)

// StartCleanupRoutine launches a background goroutine that periodically evicts
// stale entries from in-memory caches: the rate limiter map, the poster cache,
// and the SequentialMode map.
func StartCleanupRoutine() {
	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()
		for range ticker.C {
			sweepRateLimiters()
			sweepPosterCache()
			sweepSequentialMode()
		}
	}()
	slog.Info("cache cleanup routine started", "interval", cleanupInterval.String())
}

func sweepRateLimiters() {
	AuthRateLimiter.sweepOlderThan(rateLimiterMaxAge)
}

func sweepPosterCache() {
	cutoff := time.Now().Add(-posterCacheMaxAge)
	posterCache.Range(func(key, value any) bool {
		entry, ok := value.(*posterCacheEntry)
		if ok && entry.cachedAt.Before(cutoff) {
			posterCache.Delete(key)
		}
		return true
	})
}

func sweepSequentialMode() {
	Configmu.Lock()
	tc := Engine.Torc
	Configmu.Unlock()
	if tc == nil {
		return
	}
	activeTorrents := make(map[string]bool)
	for _, t := range tc.Torrents() {
		activeTorrents[t.InfoHash().HexString()] = true
	}
	SequentialMode.Range(func(key, _ any) bool {
		ih, ok := key.(string)
		if ok && !activeTorrents[ih] {
			SequentialMode.Delete(key)
		}
		return true
	})
}
