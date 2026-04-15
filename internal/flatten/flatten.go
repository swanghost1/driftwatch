// Package flatten converts nested drift results into a flat key-value
// representation suitable for tabular output or downstream processing.
package flatten

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// Row is a single flattened record derived from a drift result.
type Row struct {
	Service  string
	Field    string
	Wanted   string
	Actual   string
	Drifted  bool
}

// Result mirrors the minimal drift result shape used across the project.
type Result struct {
	Service string
	Field   string
	Wanted  string
	Actual  string
	Drifted bool
}

// Apply converts a slice of Results into a flat slice of Rows.
// Each result maps 1-to-1 with a Row; the function exists so that
// downstream consumers receive a uniform, sorted representation.
func Apply(results []Result) []Row {
	rows := make([]Row, 0, len(results))
	for _, r := range results {
		rows = append(rows, Row{
			Service: r.Service,
			Field:   r.Field,
			Wanted:  r.Wanted,
			Actual:  r.Actual,
			Drifted: r.Drifted,
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Service != rows[j].Service {
			return rows[i].Service < rows[j].Service
		}
		return rows[i].Field < rows[j].Field
	})
	return rows
}

// Write renders the flattened rows as a tab-separated table to w.
func Write(w io.Writer, rows []Row) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tFIELD\tWANTED\tACTUAL\tSTATUS")
	fmt.Fprintln(tw, "-------\t-----\t------\t------\t------")
	for _, r := range rows {
		status := "OK"
		if r.Drifted {
			status = "DRIFT"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n",
			r.Service, r.Field, r.Wanted, r.Actual, status)
	}
	return tw.Flush()
}
