// Package merge combines multiple drift result slices into a single
// deduplicated set, resolving conflicts by preferring drifted results
// over clean ones when the same service+field pair appears in more than
// one input slice.
package merge

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/example/driftwatch/internal/drift"
)

// key uniquely identifies a result by service name and field.
func key(r drift.Result) string {
	return fmt.Sprintf("%s\x00%s", r.Service, r.Field)
}

// Apply merges one or more result slices into a single slice.
// When the same service+field pair appears in multiple inputs, the
// drifted result wins; if all copies are clean the first is kept.
// The returned slice is sorted by service then field.
func Apply(slices ...[]drift.Result) []drift.Result {
	index := make(map[string]drift.Result)

	for _, results := range slices {
		for _, r := range results {
			k := key(r)
			existing, seen := index[k]
			if !seen {
				index[k] = r
				continue
			}
			// Prefer drifted over clean.
			if r.Drifted && !existing.Drifted {
				index[k] = r
			}
		}
	}

	out := make([]drift.Result, 0, len(index))
	for _, r := range index {
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

// Write prints a summary of the merged results to w.
func Write(w io.Writer, results []drift.Result) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tFIELD\tDRIFTED\tWANT\tGOT")
	for _, r := range results {
		drifted := "no"
		if r.Drifted {
			drifted = "yes"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", r.Service, r.Field, drifted, r.Want, r.Got)
	}
	tw.Flush()
}
