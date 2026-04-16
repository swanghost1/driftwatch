// Package partition splits drift results into named buckets based on a
// user-supplied key function, making it easy to process sub-groups independently.
package partition

import (
	"fmt"
	"io"
	"sort"

	"github.com/driftwatch/internal/drift"
)

// Partition holds a named slice of drift results.
type Partition struct {
	Key     string
	Results []drift.Result
}

// Options controls how results are partitioned.
type Options struct {
	// By is the partitioning dimension: "service", "field", or "status".
	By string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{By: "service"}
}

// Apply partitions results according to opts.By.
// Unknown dimensions fall back to a single partition keyed "all".
func Apply(results []drift.Result, opts Options) []Partition {
	if len(results) == 0 {
		return nil
	}

	buckets := map[string][]drift.Result{}
	for _, r := range results {
		k := keyFor(r, opts.By)
		buckets[k] = append(buckets[k], r)
	}

	keys := make([]string, 0, len(buckets))
	for k := range buckets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := make([]Partition, 0, len(keys))
	for _, k := range keys {
		out = append(out, Partition{Key: k, Results: buckets[k]})
	}
	return out
}

func keyFor(r drift.Result, by string) string {
	switch by {
	case "field":
		return r.Field
	case "status":
		if r.Drifted {
			return "drifted"
		}
		return "ok"
	default:
		return r.Service
	}
}

// Write prints a human-readable summary of partitions to w.
func Write(w io.Writer, partitions []Partition) {
	if len(partitions) == 0 {
		fmt.Fprintln(w, "no partitions")
		return
	}
	for _, p := range partitions {
		drifted := 0
		for _, r := range p.Results {
			if r.Drifted {
				drifted++
			}
		}
		fmt.Fprintf(w, "[%s] total=%d drifted=%d\n", p.Key, len(p.Results), drifted)
	}
}
