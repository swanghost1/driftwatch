// Package sample provides periodic sampling of drift results,
// recording a representative subset of runs for trend analysis
// without retaining every execution.
package sample

import (
	"math/rand"
	"time"

	"github.com/example/driftwatch/internal/drift"
)

// Options controls how sampling is applied.
type Options struct {
	// Rate is the probability [0.0, 1.0] that a run is sampled.
	// 1.0 means always sample; 0.0 means never sample.
	Rate float64

	// AlwaysSampleDrifted ensures runs with at least one drifted
	// result are always recorded, regardless of Rate.
	AlwaysSampleDrifted bool
}

// DefaultOptions returns a sensible default sampling configuration.
func DefaultOptions() Options {
	return Options{
		Rate:                0.5,
		AlwaysSampleDrifted: true,
	}
}

// ShouldRecord reports whether the current run should be persisted
// based on the sampling options and the results of the run.
func ShouldRecord(results []drift.Result, opts Options) bool {
	if opts.Rate <= 0.0 {
		return false
	}
	if opts.Rate >= 1.0 {
		return true
	}
	if opts.AlwaysSampleDrifted && anyDrifted(results) {
		return true
	}
	//nolint:gosec // non-cryptographic random is acceptable for sampling
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Float64() < opts.Rate
}

// Apply filters results down to a random sample of size n.
// If n <= 0 or n >= len(results) the original slice is returned unchanged.
// Drifted results are always preferred when trimming.
func Apply(results []drift.Result, n int) []drift.Result {
	if n <= 0 || n >= len(results) {
		return results
	}
	drifted := make([]drift.Result, 0, len(results))
	clean := make([]drift.Result, 0, len(results))
	for _, r := range results {
		if r.Drifted {
			drifted = append(drifted, r)
		} else {
			clean = append(clean, r)
		}
	}
	out := make([]drift.Result, 0, n)
	out = append(out, drifted...)
	for _, r := range clean {
		if len(out) >= n {
			break
		}
		out = append(out, r)
	}
	return out
}

func anyDrifted(results []drift.Result) bool {
	for _, r := range results {
		if r.Drifted {
			return true
		}
	}
	return false
}
