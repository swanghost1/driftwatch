package compare

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
)

// WriteText writes a human-readable summary of changes to w.
func WriteText(w io.Writer, changes []Change) error {
	if len(changes) == 0 {
		_, err := fmt.Fprintln(w, "no changes detected between snapshots")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tFIELD\tEXPECTED\tACTUAL\tCHANGE")

	for _, c := range changes {
		r := c.Result
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n",
			r.Service, r.Field, r.Expected, r.Actual, c.Kind)
	}

	return tw.Flush()
}

// WriteJSON writes changes as a JSON array to w.
func WriteJSON(w io.Writer, changes []Change) error {
	type row struct {
		Service  string     `json:"service"`
		Field    string     `json:"field"`
		Expected string     `json:"expected"`
		Actual   string     `json:"actual"`
		Change   ChangeKind `json:"change"`
	}

	rows := make([]row, 0, len(changes))
	for _, c := range changes {
		rows = append(rows, row{
			Service:  c.Result.Service,
			Field:    c.Result.Field,
			Expected: c.Result.Expected,
			Actual:   c.Result.Actual,
			Change:   c.Kind,
		})
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}
