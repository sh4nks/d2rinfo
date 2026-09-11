package server

import (
	"time"

	otter "github.com/maypok86/otter/v2"
)

// FixedIntervalExpiry is a concrete implementation of otter.ExpiryCalculator.
type FixedIntervalExpiry struct{}

// ExpireAfterCreate calculates a TTL to expire at the next boundary:
// - If created between :00 and :10, expires at :10.
// - If created between :10 and :30, expires at :30.
// - If created between :30 and :40, expires at :40.
// - If created between :40 and :00, expires at :00 of the next hour.
// While the next terror zone is not predicted yet, it expires at
// next_available_time_utc instead if that comes first.
func (e FixedIntervalExpiry) ExpireAfterCreate(entry otter.Entry[string, any]) time.Duration {
	return expiryFor(entry.Value)
}

func expiryFor(value any) time.Duration {
	now := time.Now()
	minute := now.Minute()
	var nextBoundary time.Time

	if minute < 10 {
		nextBoundary = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 10, 0, 0, now.Location())
	} else if minute < 30 {
		nextBoundary = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 30, 0, 0, now.Location())
	} else if minute < 40 {
		nextBoundary = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 40, 0, 0, now.Location())
	} else {
		nextBoundary = time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())
	}

	// A prediction that is still missing after next_available_time_utc falls
	// back to the fixed boundary, so a late d2emu.com is not polled nonstop.
	if predictionAt, ok := pendingPredictionTime(value); ok && predictionAt.After(now) && predictionAt.Before(nextBoundary) {
		nextBoundary = predictionAt
	}

	ttl := nextBoundary.Sub(now)
	if ttl <= 0 {
		return time.Minute
	}
	return ttl
}

// pendingPredictionTime returns next_available_time_utc if the cached TZ data
// has no prediction for the next terror zone yet.
func pendingPredictionTime(value any) (time.Time, bool) {
	data, ok := value.(map[string]any)
	if !ok {
		return time.Time{}, false
	}
	tz, ok := data["tz"].(map[string]any)
	if !ok {
		return time.Time{}, false
	}
	if next, _ := tz["next"].([]any); len(next) > 0 {
		return time.Time{}, false
	}
	secs, ok := tz["next_available_time_utc"].(float64)
	if !ok || secs <= 0 {
		return time.Time{}, false
	}
	return time.Unix(int64(secs), 0), true
}

// ExpireAfterRead returns the remaining duration until the entry's original expiration time.
func (e FixedIntervalExpiry) ExpireAfterRead(entry otter.Entry[string, any]) time.Duration {
	// To preserve the expiration time, we calculate the remaining duration
	// from the current time until the absolute expiration time using time.Until.
	return time.Until(entry.ExpiresAt())
}

// ExpireAfterUpdate recalculates the TTL as if the entry were new.
func (e FixedIntervalExpiry) ExpireAfterUpdate(entry otter.Entry[string, any], value any) time.Duration {
	return expiryFor(value)
}
