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

// FIXME: upper bound needs to be greater than lower bound
// FIXME: explicit that only integer values are supported
// ParseLatency parses a latency string and returns a Latency struct.
// The latency string can be in the following formats:
// 1. {value}{time unit identifier}
// 2. {value}{time unit identifier}-{value}{time unit identifier}
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

// Wait waits for the duration held by the latency.
//
// If the latency has upper and lower bounds, it waits for a random duration
// between the lower and upper bounds. Otherwise it waits for the lower bound.
func (l Latency) Wait() time.Duration {
	// If the latency has no bounds, wait for the specified duration
	if !l.HasBounds() {
		time.Sleep(l.LowerBound)
		return l.LowerBound
	}

	waitTime := time.Duration(rand.Int63n(int64(l.UpperBound-l.LowerBound))) + l.LowerBound

	// Wait for a random duration between the lower and upper bounds
	time.Sleep(waitTime)

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
