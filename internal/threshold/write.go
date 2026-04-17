package threshold

import (
	"encoding/json"
	"fmt"
	"io"
)

type jsonResult struct {
	Total     int     `json:"total"`
	Drifted   int     `json:"drifted"`
	DriftRate float64 `json:"drift_rate"`
	Breached  bool    `json:"breached"`
	Reason    string  `json:"reason,omitempty"`
}

// WriteJSON encodes the threshold result as JSON to w.
func WriteJSON(w io.Writer, r Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(jsonResult{
		Total:     r.Total,
		Drifted:   r.Drifted,
		DriftRate: r.DriftRate,
		Breached:  r.Breached,
		Reason:    r.Reason,
	})
}

// WriteText writes a compact single-line summary to w.
func WriteText(w io.Writer, r Result) {
	status := "ok"
	if r.Breached {
		status = "breached"
	}
	fmt.Fprintf(w, "threshold=%s drifted=%d/%d rate=%.2f", status, r.Drifted, r.Total, r.DriftRate)
	if r.Breached {
		fmt.Fprintf(w, " reason=%q", r.Reason)
	}
	fmt.Fprintln(w)
}
