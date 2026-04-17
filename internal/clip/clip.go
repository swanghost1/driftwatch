// Package clip limits the number of drift results retained per field across
// all services. This is useful when a single misconfigured field produces
// noise across many services and you want a representative sample only.
package clip

import (
	"fmt"
	"io"
	"sort"

	"github.com/example/driftwatch/internal/drift"
)

// DefaultOptions returns a sensible default: keep up to 5 results per field.
func DefaultOptions() Options {
	return Options{MaxPerField: 5}
}

// Options controls how clipping is applied.
type Options struct {
	// MaxPerField is the maximum number of results to keep for each distinct
	// field name. Zero means no limit.
	MaxPerField int
}

// Apply returns a new slice where each field name appears at most
// opts.MaxPerField times. Order is preserved; excess entries are dropped.
func Apply(results []drift.Result, opts Options) []drift.Result {
	if opts.MaxPerField == 0 {
		return results
	}
	counts := make(map[string]int)
	out := results[:0:0]
	for _, r := range results {
		if counts[r.Field] < opts.MaxPerField {
			out = append(out, r)
			counts[r.Field]++
		}
	}
	return out
}

// Write prints a summary of how many results were kept per field.
func Write(w io.Writer, results []drift.Result) {
	counts := make(map[string]int)
	for _, r := range results {
		counts[r.Field]++
	}
	fields := make([]string, 0, len(counts))
	for f := range counts {
		fields = append(fields, f)
	}
	sort.Strings(fields)
	fmt.Fprintln(w, "FIELD                            COUNT")
	fmt.Fprintln(w, "-------------------------------- -----")
	for _, f := range fields {
		fmt.Fprintf(w, "%-32s %5d\n", f, counts[f])
	}
}
