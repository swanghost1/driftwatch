// Package reorder sorts drift results according to a configurable priority list.
// Results are ordered by a sequence of named criteria; ties fall through to the
// next criterion in the list.
package reorder

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/example/driftwatch/internal/drift"
)

// Criterion names a single ordering dimension.
type Criterion string

const (
	ByService  Criterion = "service"
	ByField    Criterion = "field"
	ByDrifted  Criterion = "drifted"
	ByExpected Criterion = "expected"
)

// Options controls how Apply orders results.
type Options struct {
	// Order is the priority list of criteria, applied left-to-right.
	// Unknown criteria are silently ignored.
	Order []Criterion
}

// DefaultOptions returns a sensible default ordering: drifted first, then
// service name, then field name.
func DefaultOptions() Options {
	return Options{Order: []Criterion{ByDrifted, ByService, ByField}}
}

// Apply returns a new slice of results sorted according to opts.
// The original slice is not modified.
func Apply(results []drift.Result, opts Options) []drift.Result {
	out := make([]drift.Result, len(results))
	copy(out, results)

	sort.SliceStable(out, func(i, j int) bool {
		for _, c := range opts.Order {
			cmp := compareBy(out[i], out[j], c)
			if cmp != 0 {
				return cmp < 0
			}
		}
		return false
	})
	return out
}

func compareBy(a, b drift.Result, c Criterion) int {
	switch c {
	case ByService:
		return strings.Compare(a.Service, b.Service)
	case ByField:
		return strings.Compare(a.Field, b.Field)
	case ByDrifted:
		// drifted (true) sorts before clean (false)
		if a.Drifted == b.Drifted {
			return 0
		}
		if a.Drifted {
			return -1
		}
		return 1
	case ByExpected:
		return strings.Compare(a.Expected, b.Expected)
	}
	return 0
}

// Write prints a human-readable summary of the ordering criteria to w.
func Write(w io.Writer, opts Options) {
	names := make([]string, len(opts.Order))
	for i, c := range opts.Order {
		names[i] = string(c)
	}
	fmt.Fprintf(w, "result order: %s\n", strings.Join(names, " → "))
}
