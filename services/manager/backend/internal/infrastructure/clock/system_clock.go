package clock

import "time"

// SystemClock provides the current time.
type SystemClock struct{}

// NewSystemClock constructs a new clock instance.
func NewSystemClock() SystemClock {
	return SystemClock{}
}

// Now returns the current UTC time with microsecond precision.
func (SystemClock) Now() time.Time {
	return time.Now().UTC().Truncate(time.Microsecond)
}
