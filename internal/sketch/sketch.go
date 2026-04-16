// Package sketch provides a lightweight summary of drift results,
// condensing multiple fields into a single per-service overview.
package sketch

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/driftwatch/internal/drift"
)

// Entry summarises drift state for a single service.
type Entry struct {
	Service    string   `json:"service"`
	Drifted    bool     `json:"drifted"`
	Fields     []string `json:"fields,omitempty"`
	DriftCount int      `json:"drift_count"`
	Total      int      `json:"total"`
}

// Compute builds one Entry per service from the provided results.
func Compute(results []drift.Result) []Entry {
	type accum struct {
		fields []string
		total  int
	}
	m := map[string]*accum{}
	for _, r := range results {
		a, ok := m[r.Service]
		if !ok {
			a = &accum{}
			m[r.Service] = a
		}
		a.total++
		if r.Drifted {
			a.fields = append(a.fields, r.Field)
		}
	}

	entries := make([]Entry, 0, len(m))
	for svc, a := range m {
		sort.Strings(a.fields)
		entries = append(entries, Entry{
			Service:    svc,
			Drifted:    len(a.fields) > 0,
			Fields:     a.fields,
			DriftCount: len(a.fields),
			Total:      a.total,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Service < entries[j].Service
	})
	return entries
}

// Write renders entries as a human-readable table to w.
func Write(w io.Writer, entries []Entry) {
	fmt.Fprintf(w, "%-30s %-8s %s\n", "SERVICE", "DRIFTED", "FIELDS")
	for _, e := range entries {
		drifted := "no"
		if e.Drifted {
			drifted = "YES"
		}
		fields := "-"
		if len(e.Fields) > 0 {
			fields = fmt.Sprintf("%v", e.Fields)
		}
		fmt.Fprintf(w, "%-30s %-8s %s\n", e.Service, drifted, fields)
	}
}

// WriteJSON encodes entries as JSON to w.
func WriteJSON(w io.Writer, entries []Entry) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}
