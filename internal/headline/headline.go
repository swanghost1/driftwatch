// Package headline produces a single-line summary of a drift detection run
// suitable for display in terminals, notifications, or log aggregators.
package headline

import (
	"fmt"
	"io"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Summary holds the condensed result of a detection run.
type Summary struct {
	Total    int
	Drifted  int
	Clean    int
	At       time.Time
}

// Build constructs a Summary from a slice of detection results.
func Build(results []drift.Result) Summary {
	s := Summary{At: time.Now()}
	for _, r := range results {
		s.Total++
		if r.Drifted {
			s.Drifted++
		} else {
			s.Clean++
		}
	}
	return s
}

// WriteText writes a human-readable one-liner to w.
func WriteText(w io.Writer, s Summary) error {
	status := "OK"
	if s.Drifted > 0 {
		status = "DRIFT DETECTED"
	}
	_, err := fmt.Fprintf(w, "[%s] %s — %d/%d services drifted (%d clean) at %s\n",
		status,
		plural(s.Drifted, "service"),
		s.Drifted,
		s.Total,
		s.Clean,
		s.At.Format(time.RFC3339),
	)
	return err
}

func plural(n int, word string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, word)
	}
	return fmt.Sprintf("%d %ss", n, word)
}
