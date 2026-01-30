package model

import (
	"time"
)

// FixedClock provides a fixed time for deterministic testing.
type FixedClock struct {
	T time.Time
}

// Now returns the fixed time.
func (c *FixedClock) Now() time.Time {
	return c.T
}

// AdvanceableClock provides a controllable clock for testing time-dependent behavior.
type AdvanceableClock struct {
	t time.Time
}

// Now returns the current clock time.
func (c *AdvanceableClock) Now() time.Time {
	return c.t
}

// Advance moves the clock forward by the given duration.
func (c *AdvanceableClock) Advance(d time.Duration) {
	c.t = c.t.Add(d)
}
