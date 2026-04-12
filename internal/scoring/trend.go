package scoring

import (
	"fmt"
	"io"
)

// TrendEntry pairs a timestamp label with a score for trend display.
type TrendEntry struct {
	Label string
	Score Score
}

// Trend holds an ordered series of score snapshots.
type Trend []TrendEntry

// Direction summarises whether the score is improving, degrading, or stable
// by comparing the first and last entries in the trend.
func (t Trend) Direction() string {
	if len(t) < 2 {
		return "stable"
	}
	first := t[0].Score.Value
	last := t[len(t)-1].Score.Value
	switch {
	case last > first:
		return "improving"
	case last < first:
		return "degrading"
	default:
		return "stable"
	}
}

// WriteTrend renders a simple ASCII table of score history to w.
func WriteTrend(w io.Writer, t Trend) {
	if len(t) == 0 {
		fmt.Fprintln(w, "No scoring history available.")
		return
	}
	fmt.Fprintf(w, "%-24s  %8s  %5s\n", "Timestamp", "Score", "Grade")
	fmt.Fprintf(w, "%-24s  %8s  %5s\n", strings.Repeat("-", 24), "--------", "-----")
	for _, e := range t {
		fmt.Fprintf(w, "%-24s  %8.2f  %5s\n", e.Label, e.Score.Value, e.Score.Grade)
	}
	fmt.Fprintf(w, "\nTrend: %s\n", t.Direction())
}
