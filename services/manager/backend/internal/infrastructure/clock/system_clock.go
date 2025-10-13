// Package clock provides time abstractions for the application.
package clock

import (
	"sync"
	"time"
)

// SystemClock exposes time retrieval with optional overrides for testing.
type SystemClock interface {
	Now() time.Time
	Set(time.Time)
	Reset()
}

type systemClock struct {
	mu     sync.RWMutex
	frozen bool
	fixed  time.Time
}

// NewSystemClock constructs a clock that returns UTC timestamps with microsecond precision.
func NewSystemClock() SystemClock {
	return &systemClock{}
}

func (c *systemClock) Now() time.Time {
	c.mu.RLock()
	frozen := c.frozen
	fixed := c.fixed
	c.mu.RUnlock()

	if frozen {
		return fixed
	}
	return time.Now().UTC().Truncate(time.Microsecond)
}

// Set forces the clock to return the provided time until Reset is called.
func (c *systemClock) Set(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.frozen = true
	c.fixed = t.UTC().Truncate(time.Microsecond)
}

// Reset restores real time progression for the clock.
func (c *systemClock) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.frozen = false
	c.fixed = time.Time{}
}
