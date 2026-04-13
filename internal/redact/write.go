package redact

import (
	"encoding/json"
	"fmt"
	"io"
)

// WriteText writes a human-readable summary of redacted fields to w.
func WriteText(w io.Writer, results []Result) {
	redacted := 0
	for _, r := range results {
		if r.Expected == redactedPlaceholder || r.Actual == redactedPlaceholder {
			redacted++
			fmt.Fprintf(w, "  [redacted] service=%-20s field=%s\n", r.Service, r.Field)
		}
	}
	if redacted == 0 {
		fmt.Fprintln(w, "  no sensitive fields redacted")
		return
	}
	fmt.Fprintf(w, "  total redacted fields: %d\n", redacted)
}

// WriteJSON encodes the results (with redactions applied) as JSON to w.
func WriteJSON(w io.Writer, results []Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}
