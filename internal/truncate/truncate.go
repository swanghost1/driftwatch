// Package truncate provides utilities for truncating drift result sets
// to a maximum number of entries, with optional prioritisation by severity.
package truncate

import (
	"fmt"
	"io"
	"sort"

	"github.com/example/driftwatch/internal/drift"
)

// Options controls how truncation is applied.
type Options struct {
	// MaxResults is the maximum number of results to retain.
	// A value of zero means no truncation.
	MaxResults int

	// PrioritiseDrifted causes drifted results to be kept over clean ones
	// when truncation occurs.
	PrioritiseDrifted bool
}

// Apply returns a (possibly shortened) slice of results according to opts.
// When PrioritiseDrifted is true, drifted entries are sorted to the front
// before the cap is applied, preserving their original relative order.
// The original slice is never modified.
func Apply(results []drift.Result, opts Options) ([]drift.Result, int) {
	if opts.MaxResults <= 0 || len(results) <= opts.MaxResults {
		return results, 0
	}

	working := make([]drift.Result, len(results))
	copy(working, results)

	if opts.PrioritiseDrifted {
		sort.SliceStable(working, func(i, j int) bool {
			// drifted (HasDrift==true) sorts before clean
			return working[i].HasDrift && !working[j].HasDrift
		})
	}

	dropped := len(working) - opts.MaxResults
	return working[:opts.MaxResults], dropped
}

// Write prints a short summary of the truncation to w.
func Write(w io.Writer, kept, dropped int) {
	if dropped <= 0 {
		fmt.Fprintf(w, "truncate: all %d results retained\n", kept)
		return
	}
	fmt.Fprintf(w, "truncate: showing %d of %d results (%d dropped)\n",
		kept, kept+dropped, dropped)
}
