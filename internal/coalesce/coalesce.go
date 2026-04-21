// Package coalesce merges multiple drift results for the same service+field
// into a single representative result, keeping the most severe or most recent.
package coalesce

import (
	"fmt"
	"io"
	"sort"

	"github.com/driftwatch/internal/drift"
)

// Options controls how coalescing is performed.
type Options struct {
	// PreferDrifted keeps a drifted result over a clean one when merging.
	PreferDrifted bool
	// PreferNewest keeps the most recently detected result when scores are equal.
	PreferNewest bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		PreferDrifted: true,
		PreferNewest:  true,
	}
}

func key(r drift.Result) string {
	return r.Service + "\x00" + r.Field
}

// Apply coalesces results so that each service+field pair appears at most once.
// When duplicates exist the winner is chosen according to opts.
func Apply(results []drift.Result, opts Options) []drift.Result {
	if len(results) == 0 {
		return nil
	}

	best := make(map[string]drift.Result, len(results))

	for _, r := range results {
		k := key(r)
		existing, seen := best[k]
		if !seen {
			best[k] = r
			continue
		}
		if opts.PreferDrifted && r.Drifted && !existing.Drifted {
			best[k] = r
			continue
		}
		if opts.PreferDrifted && !r.Drifted && existing.Drifted {
			continue
		}
		if opts.PreferNewest && r.DetectedAt.After(existing.DetectedAt) {
			best[k] = r
		}
	}

	out := make([]drift.Result, 0, len(best))
	for _, r := range best {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Service != out[j].Service {
			return out[i].Service < out[j].Service
		}
		return out[i].Field < out[j].Field
	})
	return out
}

// Write prints a human-readable summary of the coalesced results to w.
func Write(w io.Writer, results []drift.Result) {
	drifted := 0
	for _, r := range results {
		if r.Drifted {
			drifted++
		}
	}
	fmt.Fprintf(w, "coalesced: %d results (%d drifted, %d clean)\n",
		len(results), drifted, len(results)-drifted)
	for _, r := range results {
		status := "OK"
		if r.Drifted {
			status = "DRIFT"
		}
		fmt.Fprintf(w, "  [%s] %s / %s\n", status, r.Service, r.Field)
	}
}
