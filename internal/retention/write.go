package retention

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// Report summarises a prune operation.
type Report struct {
	Pruned  []string
	Policy  Policy
	RanAt   time.Time
}

// Write renders a human-readable prune report to w.
func Write(w io.Writer, r Report) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintf(tw, "Retention prune — %s\n", r.RanAt.Format(time.RFC3339))
	fmt.Fprintf(tw, "Policy:\tmax-age=%v\tmax-entries=%d\n", r.Policy.MaxAge, r.Policy.MaxEntries)
	fmt.Fprintf(tw, "Deleted:\t%d file(s)\n", len(r.Pruned))

	if len(r.Pruned) > 0 {
		fmt.Fprintln(tw, "")
		fmt.Fprintln(tw, "Removed files:")
		for _, p := range r.Pruned {
			fmt.Fprintf(tw, "  - %s\n", p)
		}
	}

	return tw.Flush()
}
