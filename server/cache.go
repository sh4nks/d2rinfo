package server

import (
	"time"

	otter "github.com/maypok86/otter/v2"
)

// FixedIntervalExpiry is a concrete implementation of otter.ExpiryCalculator.
type FixedIntervalExpiry struct{}

// ExpireAfterCreate calculates a TTL to expire at the next 30-minute mark.
func (e FixedIntervalExpiry) ExpireAfterCreate(entry otter.Entry[string, any]) time.Duration {
	now := time.Now()
	minute := now.Minute()
	var nextBoundary time.Time

	if minute < 30 {
		// If current minute is before :30, the next boundary is :30 of the current hour.
		nextBoundary = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 30, 0, 0, now.Location())
	} else {
		// If current minute is :30 or after, the next boundary is :00 of the next hour.
		nextBoundary = time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())
	}

	ttl := nextBoundary.Sub(now)
	if ttl <= 0 {
		nextBoundary = nextBoundary.Add(30 * time.Minute)
		ttl = nextBoundary.Sub(now)
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
