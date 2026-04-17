// Package peak identifies the highest-drift runs across a result history,
// surfacing the worst observed drift episodes for a given service or field.
package peak

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/driftwatch/internal/drift"
)

// Entry holds the peak drift count observed in a single run.
type Entry struct {
	RunID      string    `json:"run_id"`
	DriftCount int       `json:"drift_count"`
	Results    []drift.Result `json:"results"`
}

// Options controls peak detection behaviour.
type Options struct {
	// TopN limits how many peak entries are returned. Zero means return all.
	TopN int
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{TopN: 5}
}

// Apply receives a map of runID -> results and returns the top-N runs
// ordered by drift count descending.
func Apply(runs map[string][]drift.Result, opts Options) []Entry {
	if len(runs) == 0 {
		return nil
	}

	entries := make([]Entry, 0, len(runs))
	for id, results := range runs {
		count := 0
		for _, r := range results {
			if r.Drifted {
				count++
			}
		}
		entries = append(entries, Entry{RunID: id, DriftCount: count, Results: results})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].DriftCount != entries[j].DriftCount {
			return entries[i].DriftCount > entries[j].DriftCount
		}
		return entries[i].RunID < entries[j].RunID
	})

	if opts.TopN > 0 && len(entries) > opts.TopN {
		entries = entries[:opts.TopN]
	}
	return entries
}

// Write renders peak entries as human-readable text.
func Write(w io.Writer, entries []Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "no peak drift entries found")
		return
	}
	fmt.Fprintf(w, "%-36s  %s\n", "RUN ID", "DRIFT COUNT")
	fmt.Fprintf(w, "%-36s  %s\n", "------", "-----------")
	for _, e := range entries {
		fmt.Fprintf(w, "%-36s  %d\n", e.RunID, e.DriftCount)
	}
}

// WriteJSON renders peak entries as JSON.
func WriteJSON(w io.Writer, entries []Entry) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}
