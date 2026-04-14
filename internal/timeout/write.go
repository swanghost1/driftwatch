package timeout

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Record holds a summary of a single timed detection run.
type Record struct {
	Deadline    time.Duration `json:"deadline"`
	GracePeriod time.Duration `json:"grace_period"`
	Elapsed     time.Duration `json:"elapsed"`
	TimedOut    bool          `json:"timed_out"`
}

// WriteText writes a human-readable timeout summary to w.
func WriteText(w io.Writer, r Record) error {
	status := "ok"
	if r.TimedOut {
		status = "TIMED OUT"
	}
	_, err := fmt.Fprintf(w,
		"timeout: status=%-10s elapsed=%-10s deadline=%-10s grace=%s\n",
		status,
		r.Elapsed.Round(time.Millisecond),
		r.Deadline.Round(time.Millisecond),
		r.GracePeriod.Round(time.Millisecond),
	)
	return err
}

// WriteJSON writes r as a single JSON object to w.
func WriteJSON(w io.Writer, r Record) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
