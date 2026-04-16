package jitter

import (
	"math/rand"
	"time"
)

// Options controls how jitter is applied to a base duration.
type Options struct {
	// Factor is the maximum fraction of the base duration to add as jitter.
	// e.g. 0.2 means up to ±20% of the base duration. Defaults to 0.1.
	Factor float64
	// Seed is used to initialise the random source. 0 means use a time-based seed.
	Seed int64
}

// DefaultOptions returns sensible defaults for jitter.
func DefaultOptions() Options {
	return Options{Factor: 0.1}
}

// Apply adds a random jitter to base, returning the adjusted duration.
// The jitter is uniformly distributed in [0, Factor*base].
func Apply(base time.Duration, opts Options) time.Duration {
	if opts.Factor <= 0 {
		return base
	}
	var r *rand.Rand
	if opts.Seed != 0 {
		r = rand.New(rand.NewSource(opts.Seed)) //nolint:gosec
	} else {
		r = rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	}
	max := float64(base) * opts.Factor
	offset := time.Duration(r.Float64() * max)
	return base + offset
}

// ApplyFull adds jitter in the range [-Factor*base, +Factor*base].
func ApplyFull(base time.Duration, opts Options) time.Duration {
	if opts.Factor <= 0 {
		return base
	}
	var r *rand.Rand
	if opts.Seed != 0 {
		r = rand.New(rand.NewSource(opts.Seed)) //nolint:gosec
	} else {
		r = rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	}
	max := float64(base) * opts.Factor
	offset := time.Duration((r.Float64()*2 - 1) * max)
	return base + offset
}
