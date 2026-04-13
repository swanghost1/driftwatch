package groupfilter

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/example/driftwatch/internal/drift"
)

// WriteGroupSummary writes a human-readable breakdown of drift results
// organised by group to w.
func WriteGroupSummary(w io.Writer, results []drift.Result) {
	type groupStats struct {
		total   int
		drifted int
	}

	stats := map[string]*groupStats{}
	for _, r := range results {
		g := groupOf(r.Service)
		if stats[g] == nil {
			stats[g] = &groupStats{}
		}
		stats[g].total++
		if len(r.Fields) > 0 {
			stats[g].drifted++
		}
	}

	keys := make([]string, 0, len(stats))
	for k := range stats {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintln(w, "GROUP SUMMARY")
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, g := range keys {
		s := stats[g]
		status := "OK"
		if s.drifted > 0 {
			status = "DRIFT"
		}
		fmt.Fprintf(w, "%-20s %s  (%d/%d drifted)\n", g, status, s.drifted, s.total)
	}
}
