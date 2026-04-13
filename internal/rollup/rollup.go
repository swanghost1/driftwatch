// Package rollup aggregates drift results across multiple services into
// a concise per-namespace (or per-tag) summary suitable for dashboards.
package rollup

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/example/driftwatch/internal/drift"
)

// Group holds aggregated drift statistics for a named group.
type Group struct {
	Name       string
	Total      int
	Drifted    int
	Clean      int
	DriftRate  float64
}

// ByGroup aggregates drift.Result entries by the value returned by keyFn.
// keyFn receives the service name and should return the group key (e.g. a
// namespace prefix or tag value). Results whose keyFn returns "" are placed
// in a group named "(untagged)".
func ByGroup(results []drift.Result, keyFn func(service string) string) []Group {
	counts := make(map[string][2]int) // [total, drifted]

	for _, r := range results {
		key := keyFn(r.Service)
		if key == "" {
			key = "(untagged)"
		}
		c := counts[key]
		c[0]++
		if r.Drifted {
			c[1]++
		}
		counts[key] = c
	}

	groups := make([]Group, 0, len(counts))
	for name, c := range counts {
		total, drifted := c[0], c[1]
		rate := 0.0
		if total > 0 {
			rate = float64(drifted) / float64(total) * 100
		}
		groups = append(groups, Group{
			Name:      name,
			Total:     total,
			Drifted:   drifted,
			Clean:     total - drifted,
			DriftRate: rate,
		})
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})
	return groups
}

// Write renders the rollup table to w.
func Write(w io.Writer, groups []Group) errorw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "GROUP\tTOTAL\tDRIFTED\tCLEAN\tDRIFT%")
	fmt.Fprintln(tw, "-----\t-----\t-------\t-----\t------")
	for _, g := range groups {
		fmt.Fprintf(tw, "%s\t%d\t%d\t%d\t%.1f%%\n",
			g.Name, g.Total, g.Drifted, g.Clean, g.DriftRate)
	}
	return tw.Flush()
}
