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

// Summary returns a single-line string describing the prune result,
// suitable for logging or structured output.
func (r Report) Summary() string {
	if len(r.Pruned) == 0 {
		return fmt.Sprintf("retention prune at %s: nothing removed", r.RanAt.Format(time.RFC3339))
	}
	return fmt.Sprintf("retention prune at %s: removed %d file(s)", r.RanAt.Format(time.RFC3339), len(r.Pruned))
}
