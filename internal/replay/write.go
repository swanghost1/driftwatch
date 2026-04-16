package replay

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
)

// WriteText writes a human-readable replay summary to w.
func WriteText(w io.Writer, results []Result) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tFIELD\tDRIFTED\tREPLAYED AT")
	fmt.Fprintln(tw, "-------\t-----\t-------\t-----------")
	for _, r := range results {
		drifted := "no"
		if r.Drifted {
			drifted = "YES"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			r.Service,
			r.Field,
			drifted,
			r.ReplayedAt.Format("2006-01-02T15:04:05Z"),
		)
	}
	return tw.Flush()
}

// WriteJSON writes results as a JSON array to w.
func WriteJSON(w io.Writer, results []Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}
