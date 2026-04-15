// Package ceiling enforces an upper bound on the number of drift results
// returned across all services combined, trimming by severity then service name.
package ceiling

import (
	"fmt"
	"io"
	"sort"

	"github.com/example/driftwatch/internal/drift"
)

// Options controls ceiling behaviour.
type Options struct {
	// MaxResults is the total number of results allowed. Zero means unlimited.
	MaxResults int
	// DriftedFirst ensures drifted results are retained before clean ones when
	// the ceiling is applied.
	DriftedFirst bool
}

// DefaultOptions returns a ceiling with no limit applied.
func DefaultOptions() Options {
	return Options{MaxResults: 0, DriftedFirst: true}
}

// Apply trims results to at most opts.MaxResults entries.
// When DriftedFirst is true, drifted results are sorted before clean ones so
// they survive the cut. Within each tier results are kept in their original
// order.
func Apply(results []drift.Result, opts Options) []drift.Result {
	if opts.MaxResults <= 0 || len(results) <= opts.MaxResults {
		return results
	}

	working := make([]drift.Result, len(results))
	copy(working, results)

	if opts.DriftedFirst {
		sort.SliceStable(working, func(i, j int) bool {
			if working[i].Drifted == working[j].Drifted {
				return false
			}
			return working[i].Drifted
		})
	}

	return working[:opts.MaxResults]
}

// Write prints a short summary of the ceiling operation to w.
func Write(w io.Writer, total, kept int, opts Options) {
	if opts.MaxResults <= 0 {
		fmt.Fprintf(w, "ceiling: disabled (showing all %d results)\n", total)
		return
	}
	trimmed := total - kept
	fmt.Fprintf(w, "ceiling: max=%d total=%d kept=%d trimmed=%d\n",
		opts.MaxResults, total, kept, trimmed)
}
