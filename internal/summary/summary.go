// Package summary produces a high-level run summary suitable for
// logging or downstream consumption.
package summary

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/driftwatch/driftwatch/internal/drift"
)

// Report holds aggregated counts for a single detection run.
type Report struct {
	RunAt      time.Time `json:"run_at"`
	Total      int       `json:"total"`
	Drifted    int       `json:"drifted"`
	Clean      int       `json:"clean"`
	DriftRate  float64   `json:"drift_rate"`
}

// Build constructs a Report from a slice of detection results.
func Build(results []drift.Result) Report {
	r := Report{
		RunAt: time.Now().UTC(),
		Total: len(results),
	}
	for _, res := range results {
		if res.Drifted {
			r.Drifted++
		} else {
			r.Clean++
		}
	}
	if r.Total > 0 {
		r.DriftRate = float64(r.Drifted) / float64(r.Total) * 100
	}
	return r
}

// WriteText writes a human-readable summary to w.
func WriteText(w io.Writer, r Report) {
	fmt.Fprintf(w, "Run at:     %s\n", r.RunAt.Format(time.RFC3339))
	fmt.Fprintf(w, "Total:      %d\n", r.Total)
	fmt.Fprintf(w, "Drifted:    %d\n", r.Drifted)
	fmt.Fprintf(w, "Clean:      %d\n", r.Clean)
	fmt.Fprintf(w, "Drift rate: %.1f%%\n", r.DriftRate)
}

// WriteJSON writes the summary as a JSON object to w.
func WriteJSON(w io.Writer, r Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
