// Package pivot provides utilities for pivoting drift results by a chosen
// dimension (service or field), producing a cross-tabulated summary that
// makes it easy to spot which services or fields drift most frequently.
package pivot

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/yourorg/driftwatch/internal/drift"
)

// Axis controls which dimension is used as the pivot key.
type Axis string

const (
	AxisService Axis = "service"
	AxisField   Axis = "field"
)

// Cell holds aggregated counts for a single pivot key.
type Cell struct {
	Key      string `json:"key"`
	Total    int    `json:"total"`
	Drifted  int    `json:"drifted"`
	Clean    int    `json:"clean"`
	DriftPct float64 `json:"drift_pct"`
}

// Table is the full pivot result.
type Table struct {
	Axis  Axis   `json:"axis"`
	Cells []Cell `json:"cells"`
}

// Compute builds a pivot Table from results along the given axis.
func Compute(results []drift.Result, axis Axis) Table {
	counts := make(map[string]*Cell)

	for _, r := range results {
		key := r.Service
		if axis == AxisField {
			key = r.Field
		}
		if key == "" {
			key = "(unknown)"
		}
		c, ok := counts[key]
		if !ok {
			c = &Cell{Key: key}
			counts[key] = c
		}
		c.Total++
		if r.Drifted {
			c.Drifted++
		} else {
			c.Clean++
		}
	}

	cells := make([]Cell, 0, len(counts))
	for _, c := range counts {
		if c.Total > 0 {
			c.DriftPct = float64(c.Drifted) / float64(c.Total) * 100
		}
		cells = append(cells, *c)
	}
	sort.Slice(cells, func(i, j int) bool {
		if cells[i].Drifted != cells[j].Drifted {
			return cells[i].Drifted > cells[j].Drifted
		}
		return cells[i].Key < cells[j].Key
	})

	return Table{Axis: axis, Cells: cells}
}

// Write renders the Table as a human-readable text table to w.
func Write(w io.Writer, t Table) error {
	fmt.Fprintf(w, "Pivot by %s\n", t.Axis)
	fmt.Fprintf(w, "%-30s %6s %7s %6s %8s\n", "KEY", "TOTAL", "DRIFTED", "CLEAN", "DRIFT%")
	fmt.Fprintf(w, "%s\n", repeatChar('-', 62))
	for _, c := range t.Cells {
		fmt.Fprintf(w, "%-30s %6d %7d %6d %7.1f%%\n",
			c.Key, c.Total, c.Drifted, c.Clean, c.DriftPct)
	}
	return nil
}

// WriteJSON encodes the Table as JSON to w.
func WriteJSON(w io.Writer, t Table) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(t)
}

func repeatChar(ch rune, n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = ch
	}
	return string(b)
}
