// Package topk selects the top-K most frequently drifting services
// from a set of drift results, ranked by drift count.
package topk

import (
	"fmt"
	"io"
	"sort"

	"github.com/example/driftwatch/internal/drift"
)

// DefaultOptions returns a sensible default: top 5 services.
func DefaultOptions() Options {
	return Options{K: 5}
}

// Options controls the Top-K selection.
type Options struct {
	// K is the maximum number of services to return. Zero means return all.
	K int
}

// Entry holds a service name and its drift count.
type Entry struct {
	Service    string
	DriftCount int
}

// Apply counts drift occurrences per service and returns the top-K entries
// sorted descending by drift count.
func Apply(results []drift.Result, opts Options) []Entry {
	counts := make(map[string]int)
	for _, r := range results {
		if r.Drifted {
			counts[r.Service]++
		}
	}

	entries := make([]Entry, 0, len(counts))
	for svc, n := range counts {
		entries = append(entries, Entry{Service: svc, DriftCount: n})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].DriftCount != entries[j].DriftCount {
			return entries[i].DriftCount > entries[j].DriftCount
		}
		return entries[i].Service < entries[j].Service
	})

	if opts.K > 0 && len(entries) > opts.K {
		entries = entries[:opts.K]
	}
	return entries
}

// Write renders the top-K entries as a text table to w.
func Write(w io.Writer, entries []Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "no drifted services found")
		return
	}
	fmt.Fprintf(w, "%-30s  %s\n", "SERVICE", "DRIFT COUNT")
	fmt.Fprintf(w, "%-30s  %s\n", "-------", "-----------")
	for _, e := range entries {
		fmt.Fprintf(w, "%-30s  %d\n", e.Service, e.DriftCount)
	}
}
