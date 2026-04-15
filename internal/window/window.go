package window

import (
	"time"

	"github.com/driftwatch/internal/drift"
)

// Options configures the time-window filter.
type Options struct {
	// Since discards results whose DetectedAt is before this time.
	Since time.Time
	// Until discards results whose DetectedAt is after this time.
	Until time.Time
}

// Result wraps a drift result with its detected timestamp.
type Result = drift.Result

// Apply returns only the results whose DetectedAt falls within [Since, Until].
// A zero value for Since or Until means that bound is open.
func Apply(results []Result, opts Options) []Result {
	out := make([]Result, 0, len(results))
	for _, r := range results {
		if !opts.Since.IsZero() && r.DetectedAt.Before(opts.Since) {
			continue
		}
		if !opts.Until.IsZero() && r.DetectedAt.After(opts.Until) {
			continue
		}
		out = append(out, r)
	}
	return out
}

// Last returns the results detected within the most recent duration d.
func Last(results []Result, d time.Duration) []Result {
	if d <= 0 {
		return results
	}
	return Apply(results, Options{Since: time.Now().Add(-d)})
}
