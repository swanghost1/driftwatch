// Package aggregate combines drift results across multiple services into
// a single unified summary, grouping by field and computing totals.
package aggregate

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/driftwatch/internal/drift"
)

// FieldSummary holds aggregated drift counts for a single field name.
type FieldSummary struct {
	Field    string
	Total    int
	Drifted  int
	Services []string
}

// DriftRate returns the fraction of services drifted on this field.
func (f FieldSummary) DriftRate() float64 {
	if f.Total == 0 {
		return 0
	}
	return float64(f.Drifted) / float64(f.Total)
}

// ByField aggregates a slice of drift results by field name.
// Each entry in the returned slice covers all services that reported
// that field (whether drifted or not).
func ByField(results []drift.Result) []FieldSummary {
	type entry struct {
		total    int
		drifted  int
		services []string
	}

	m := make(map[string]*entry)

	for _, r := range results {
		e, ok := m[r.Field]
		if !ok {
			e = &entry{}
			m[r.Field] = e
		}
		e.total++
		e.services = append(e.services, r.Service)
		if r.Drifted {
			e.drifted++
		}
	}

	summaries := make([]FieldSummary, 0, len(m))
	for field, e := range m {
		sort.Strings(e.services)
		summaries = append(summaries, FieldSummary{
			Field:    field,
			Total:    e.total,
			Drifted:  e.drifted,
			Services: e.services,
		})
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Field < summaries[j].Field
	})

	return summaries
}

// Write renders the field-level aggregate summary to w as a text table.
func Write(w io.Writer, summaries []FieldSummary) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "FIELD\tTOTAL\tDRIFTED\tRATE")
	for _, s := range summaries {
		fmt.Fprintf(tw, "%s\t%d\t%d\t%.0f%%\n",
			s.Field, s.Total, s.Drifted, s.DriftRate()*100)
	}
	return tw.Flush()
}
