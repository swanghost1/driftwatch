// Package fanin merges multiple slices of drift results into a single
// deduplicated, ordered slice. It is useful when results arrive from
// parallel detection passes and need to be consolidated before reporting.
package fanin

import (
	"sort"

	"github.com/driftwatch/internal/drift"
)

// Options controls how the merge is performed.
type Options struct {
	// DriftedFirst ensures drifted results appear before clean ones.
	DriftedFirst bool
	// DeduplicateByKey removes duplicate (service, field) pairs, keeping the
	// first occurrence (which should be the most authoritative source).
	DeduplicateByKey bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		DriftedFirst:     true,
		DeduplicateByKey: true,
	}
}

// Apply merges all provided result slices into one according to opts.
func Apply(sources [][]drift.Result, opts Options) []drift.Result {
	var merged []drift.Result
	seen := make(map[string]struct{})

	for _, src := range sources {
		for _, r := range src {
			if opts.DeduplicateByKey {
				k := r.Service + "\x00" + r.Field
				if _, exists := seen[k]; exists {
					continue
				}
				seen[k] = struct{}{}
			}
			merged = append(merged, r)
		}
	}

	if opts.DriftedFirst {
		sort.SliceStable(merged, func(i, j int) bool {
			if merged[i].Drifted == merged[j].Drifted {
				return false
			}
			return merged[i].Drifted
		})
	}

	return merged
}
