package annotate

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/example/driftwatch/internal/drift"
)

// WriteText writes a human-readable summary of annotations present across all
// results to w.
func WriteText(w io.Writer, results []drift.Result) {
	type row struct {
		service string
		key     string
		value   string
	}
	var rows []row
	for _, r := range results {
		keys := make([]string, 0, len(r.Annotations))
		for k := range r.Annotations {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			rows = append(rows, row{r.Service, k, r.Annotations[k]})
		}
	}
	if len(rows) == 0 {
		fmt.Fprintln(w, "annotations: none")
		return
	}
	fmt.Fprintln(w, "ANNOTATIONS")
	fmt.Fprintf(w, "%-30s %-20s %s\n", "SERVICE", "KEY", "VALUE")
	for _, row := range rows {
		fmt.Fprintf(w, "%-30s %-20s %s\n", row.service, row.key, row.value)
	}
}

// WriteJSON encodes the annotation map for every result as a JSON array to w.
func WriteJSON(w io.Writer, results []drift.Result) error {
	type entry struct {
		Service     string            `json:"service"`
		Annotations map[string]string `json:"annotations"`
	}
	entries := make([]entry, 0, len(results))
	for _, r := range results {
		entries = append(entries, entry{Service: r.Service, Annotations: r.Annotations})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}
