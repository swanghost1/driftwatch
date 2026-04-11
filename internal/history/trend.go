package history

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// RunSummary holds aggregated statistics for a single history entry.
type RunSummary struct {
	RunAt      time.Time
	Total      int
	Drifted    int
	Clean      int
	DriftRate  float64
}

// Summarise converts a list of history entries into run summaries.
func Summarise(entries []Entry) []RunSummary {
	summaries := make([]RunSummary, 0, len(entries))
	for _, e := range entries {
		s := RunSummary{
			RunAt: e.RunAt,
			Total: len(e.Results),
		}
		for _, r := range e.Results {
			if r.Status == "drifted" {
				s.Drifted++
			} else {
				s.Clean++
			}
		}
		if s.Total > 0 {
			s.DriftRate = float64(s.Drifted) / float64(s.Total) * 100
		}
		summaries = append(summaries, s)
	}
	return summaries
}

// WriteTrend writes a human-readable trend table to w.
func WriteTrend(w io.Writer, summaries []RunSummary) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "RUN AT\tTOTAL\tDRIFTED\tCLEAN\tDRIFT RATE")
	for _, s := range summaries {
		fmt.Fprintf(tw, "%s\t%d\t%d\t%d\t%.1f%%\n",
			s.RunAt.Format(time.RFC3339),
			s.Total,
			s.Drifted,
			s.Clean,
			s.DriftRate,
		)
	}
	return tw.Flush()
}
