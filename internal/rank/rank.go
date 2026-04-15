// Package rank orders drift results by configurable priority criteria,
// allowing operators to surface the most important drifts first.
package rank

import (
	"io"
	"fmt"
	"sort"
	"strings"

	"github.com/yourorg/driftwatch/internal/drift"
)

// Criterion defines how results should be ordered.
type Criterion string

const (
	BySeverity  Criterion = "severity"
	ByService   Criterion = "service"
	ByField     Criterion = "field"
	ByDriftOnly Criterion = "drifted_first"
)

// Options controls ranking behaviour.
type Options struct {
	// Criteria is an ordered list of sort keys applied left-to-right.
	Criteria []Criterion
}

// DefaultOptions returns a sensible default ranking.
func DefaultOptions() Options {
	return Options{
		Criteria: []Criterion{ByDriftOnly, BySeverity, ByService},
	}
}

// Apply returns a new slice of results sorted according to opts.
// The original slice is not modified.
func Apply(results []drift.Result, opts Options) []drift.Result {
	out := make([]drift.Result, len(results))
	copy(out, results)

	sort.SliceStable(out, func(i, j int) bool {
		for _, c := range opts.Criteria {
			cmp := compare(out[i], out[j], c)
			if cmp != 0 {
				return cmp < 0
			}
		}
		return false
	})
	return out
}

func compare(a, b drift.Result, c Criterion) int {
	switch c {
	case ByDriftOnly:
		// drifted results sort before clean
		if a.Drifted == b.Drifted {
			return 0
		}
		if a.Drifted {
			return -1
		}
		return 1
	case BySeverity:
		return strings.Compare(
			normaliseSeverity(a.Severity),
			normaliseSeverity(b.Severity),
		)
	case ByService:
		return strings.Compare(a.Service, b.Service)
	case ByField:
		return strings.Compare(a.Field, b.Field)
	}
	return 0
}

// normaliseSeverity maps severity labels to a sortable key (lower = higher priority).
func normaliseSeverity(s string) string {
	switch strings.ToLower(s) {
	case "critical":
		return "0"
	case "high":
		return "1"
	case "medium":
		return "2"
	case "low":
		return "3"
	}
	return "4"
}

// Write prints a ranked summary table to w.
func Write(w io.Writer, results []drift.Result) {
	fmt.Fprintf(w, "%-30s %-20s %-10s %s\n", "SERVICE", "FIELD", "SEVERITY", "STATUS")
	fmt.Fprintln(w, strings.Repeat("-", 75))
	for _, r := range results {
		status := "ok"
		if r.Drifted {
			status = "DRIFT"
		}
		sev := r.Severity
		if sev == "" {
			sev = "-"
		}
		fmt.Fprintf(w, "%-30s %-20s %-10s %s\n", r.Service, r.Field, sev, status)
	}
}
