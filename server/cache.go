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
func (e FixedIntervalExpiry) ExpireAfterCreate(entry otter.Entry[string, any]) time.Duration {
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

	ttl := nextBoundary.Sub(now)
	if ttl <= 0 {
		return time.Minute
	}
	return ttl
}

// ExpireAfterRead returns the remaining duration until the entry's original expiration time.
func (e FixedIntervalExpiry) ExpireAfterRead(entry otter.Entry[string, any]) time.Duration {
	// To preserve the expiration time, we calculate the remaining duration
	// from the current time until the absolute expiration time using time.Until.
	return time.Until(entry.ExpiresAt())
}

// ExpireAfterUpdate recalculates the TTL as if the entry were new.
func (e FixedIntervalExpiry) ExpireAfterUpdate(entry otter.Entry[string, any], value any) time.Duration {
	return e.ExpireAfterCreate(entry)
}
