package reorder

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/example/driftwatch/internal/drift"
)

// WriteText writes an ordered result list in a human-readable table to w.
func WriteText(w io.Writer, results []drift.Result) {
	fmt.Fprintf(w, "%-30s %-15s %-8s %s\n", "SERVICE", "FIELD", "DRIFTED", "EXPECTED")
	fmt.Fprintf(w, "%s\n", repeatChar('-', 72))
	for _, r := range results {
		drifted := "no"
		if r.Drifted {
			drifted = "YES"
		}
		fmt.Fprintf(w, "%-30s %-15s %-8s %s\n", r.Service, r.Field, drifted, r.Expected)
	}
}

// WriteJSON writes the ordered results as a JSON array to w.
func WriteJSON(w io.Writer, results []drift.Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}

func repeatChar(ch rune, n int) string {
	out := make([]rune, n)
	for i := range out {
		out[i] = ch
	}
	return string(out)
}
