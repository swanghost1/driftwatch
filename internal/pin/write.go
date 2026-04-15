package pin

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// WriteText writes a human-readable summary of active pin rules to w.
func WriteText(w io.Writer, s *Store) error {
	active := activeRules(s, time.Now())
	if len(active) == 0 {
		_, err := fmt.Fprintln(w, "No active pin rules.")
		return err
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tFIELD\tREASON\tEXPIRES")
	for _, r := range active {
		expires := "never"
		if r.ExpiresAt != nil {
			expires = r.ExpiresAt.Format(time.RFC3339)
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", r.Service, fieldOrAll(r.Field), r.Reason, expires)
	}
	return tw.Flush()
}

// WriteJSON writes the active pin rules as a JSON array to w.
func WriteJSON(w io.Writer, s *Store) error {
	active := activeRules(s, time.Now())
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(active)
}

func activeRules(s *Store, now time.Time) []Rule {
	out := make([]Rule, 0, len(s.Rules))
	for _, r := range s.Rules {
		if r.ExpiresAt != nil && now.After(*r.ExpiresAt) {
			continue
		}
		out = append(out, r)
	}
	return out
}

func fieldOrAll(f string) string {
	if f == "" {
		return "(all)"
	}
	return f
}
