//nolint:revive
package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Latency represents a latency duration.
type Latency struct {
	LowerBound time.Duration
	UpperBound time.Duration
}

// ParseLatency parses a duration string and returns a Latency struct.
//
// The duration string should be in the format "lower-upper", where "lower" and "upper" are time durations.
//
// If the duration string is invalid or the upper bound is less than the lower bound, an error is returned.
func ParseLatency(duration string) (latency Latency, err error) {
	// Split the duration into the lower and upper bounds
	bounds := strings.Split(duration, "-")

	// Parse the lower bound
	latency.LowerBound, err = time.ParseDuration(bounds[0])
	if err != nil {
		err = fmt.Errorf("failed parsing duration's lower bound: %w", err)
		return
	}

	// If there is no upper bound, return
	if len(bounds) == 1 {
		return
	}

	// Parse the upper bound
	latency.UpperBound, err = time.ParseDuration(bounds[1])
	if err != nil {
		err = fmt.Errorf("failed parsing duration's lower bound: %w", err)
		return
	}

	return
}

// Validate checks if the Latency struct satisfies the defined constraints.
func (l Latency) Validate() error {
	if l.LowerBound < 0 {
		return fmt.Errorf("lower bound is negative: %w", ErrNegativeBound)
	}

	if l.UpperBound < 0 {
		return fmt.Errorf("upper bound is negative: %w", ErrNegativeBound)
	}

	if l.UpperBound > 0 && l.LowerBound > l.UpperBound {
		return ErrUpperBoundGreaterThanLowerBound
	}

	return nil
}

// Wait waits for the duration held by the latency.
//
// If the latency has upper and lower bounds, it waits for a random duration
// between the lower and upper bounds. Otherwise it waits for the lower bound.
//
//nolint:gosec
func (l Latency) Wait() time.Duration {
	// If the latency has no bounds, wait for the specified duration
	if !l.HasBounds() {
		time.Sleep(l.LowerBound)
		return l.LowerBound
	}

	waitTime := time.Duration(rand.Int63n(int64(l.UpperBound-l.LowerBound))) + l.LowerBound

	// Wait for a random duration between the lower and upper bounds
	<-time.After(waitTime)

	return waitTime
}

// String returns the string representation of the latency
func (l Latency) String() string {
	if !l.HasBounds() {
		return l.LowerBound.String()
	}

	return fmt.Sprintf("%s-%s", l.LowerBound, l.UpperBound)
}

// HasBounds returns true if the latency has upper and lower bounds
// and false otherwise.
func (l *Latency) HasBounds() bool {
	// Neither bound is set
	if l.LowerBound == 0 && l.UpperBound == 0 {
		return false
	}

	if l.LowerBound != 0 && l.UpperBound != 0 {
		return true
	}

	return false
}
