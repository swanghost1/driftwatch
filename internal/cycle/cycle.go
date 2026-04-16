// Package cycle detects repeated drift patterns across consecutive runs,
// flagging services whose drift has persisted without resolution.
package cycle

import (
	"fmt"
	"io"
	"sort"

	"github.com/driftwatch/internal/drift"
)

// Options controls cycle detection behaviour.
type Options struct {
	// MinRuns is the minimum number of consecutive drifted runs required
	// before a result is considered a persistent cycle. Default: 3.
	MinRuns int
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{MinRuns: 3}
}

// Result describes a service/field pair that has drifted persistently.
type Result struct {
	Service    string
	Field      string
	Consecutive int
}

// Detect inspects a slice of historical run results (oldest first) and
// returns every service+field pair that has been drifted in at least
// opts.MinRuns consecutive trailing runs.
func Detect(history [][]drift.Result, opts Options) []Result {
	if opts.MinRuns <= 0 {
		opts.MinRuns = DefaultOptions().MinRuns
	}
	if len(history) < opts.MinRuns {
		return nil
	}

	// Count trailing consecutive drifts per key.
	type key struct{ service, field string }
	counts := map[key]int{}

	for i := len(history) - 1; i >= 0; i-- {
		seen := map[key]bool{}
		for _, r := range history[i] {
			if !r.Drifted {
				continue
			}
			k := key{r.Service, r.Field}
			seen[k] = true
			counts[k]++
		}
		// Stop counting keys that didn't appear as drifted in this run.
		for k := range counts {
			if !seen[k] {
				delete(counts, k)
			}
		}
	}

	var out []Result
	for k, n := range counts {
		if n >= opts.MinRuns {
			out = append(out, Result{Service: k.service, Field: k.field, Consecutive: n})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Service != out[j].Service {
			return out[i].Service < out[j].Service
		}
		return out[i].Field < out[j].Field
	})
	return out
}

// Write prints cycle results in a human-readable table to w.
func Write(w io.Writer, results []Result) {
	if len(results) == 0 {
		fmt.Fprintln(w, "no persistent drift cycles detected")
		return
	}
	fmt.Fprintf(w, "%-30s %-20s %s\n", "SERVICE", "FIELD", "CONSECUTIVE RUNS")
	for _, r := range results {
		fmt.Fprintf(w, "%-30s %-20s %d\n", r.Service, r.Field, r.Consecutive)
	}
}
