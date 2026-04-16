// Package offset provides utilities for shifting drift result timestamps
// by a fixed duration. This is useful when replaying historical data or
// adjusting for timezone differences in stored snapshots.
package offset

import (
	"fmt"
	"io"
	"time"

	"driftwatch/internal/drift"
)

// Options controls how the offset is applied.
type Options struct {
	// Shift is the duration added to each result's DetectedAt timestamp.
	// Negative values shift backwards in time.
	Shift time.Duration

	// SkipClean, when true, only adjusts drifted results.
	SkipClean bool
}

// DefaultOptions returns an Options with no shift applied.
func DefaultOptions() Options {
	return Options{}
}

// Apply returns a copy of results with DetectedAt shifted by opts.Shift.
// If opts.Shift is zero the original slice is returned unchanged.
func Apply(results []drift.Result, opts Options) []drift.Result {
	if opts.Shift == 0 {
		return results
	}
	out := make([]drift.Result, len(results))
	for i, r := range results {
		if opts.SkipClean && !r.Drifted {
			out[i] = r
			continue
		}
		r.DetectedAt = r.DetectedAt.Add(opts.Shift)
		out[i] = r
	}
	return out
}

// Write writes a human-readable summary of the applied offset to w.
func Write(w io.Writer, opts Options, count int) {
	direction := "forward"
	shift := opts.Shift
	if shift < 0 {
		direction = "backward"
		shift = -shift
	}
	fmt.Fprintf(w, "offset: shifted %d result(s) %s by %s\n", count, direction, shift)
	if opts.SkipClean {
		fmt.Fprintln(w, "offset: clean results were not adjusted")
	}
}
