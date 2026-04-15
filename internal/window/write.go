package window

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Summary describes the active window bounds.
type Summary struct {
	Since    *time.Time `json:"since,omitempty"`
	Until    *time.Time `json:"until,omitempty"`
	Matched  int        `json:"matched"`
	Filtered int        `json:"filtered"`
}

// NewSummary builds a Summary given the original and filtered result slices.
func NewSummary(original, filtered []Result, opts Options) Summary {
	s := Summary{
		Matched:  len(filtered),
		Filtered: len(original) - len(filtered),
	}
	if !opts.Since.IsZero() {
		t := opts.Since
		s.Since = &t
	}
	if !opts.Until.IsZero() {
		t := opts.Until
		s.Until = &t
	}
	return s
}

// WriteText writes a human-readable summary to w.
func WriteText(w io.Writer, s Summary) error {
	fmt.Fprintln(w, "=== Window Filter ===")
	if s.Since != nil {
		fmt.Fprintf(w, "  Since   : %s\n", s.Since.Format(time.RFC3339))
	} else {
		fmt.Fprintln(w, "  Since   : (open)")
	}
	if s.Until != nil {
		fmt.Fprintf(w, "  Until   : %s\n", s.Until.Format(time.RFC3339))
	} else {
		fmt.Fprintln(w, "  Until   : (open)")
	}
	fmt.Fprintf(w, "  Matched : %d\n", s.Matched)
	fmt.Fprintf(w, "  Filtered: %d\n", s.Filtered)
	return nil
}

// WriteJSON writes the summary as JSON to w.
func WriteJSON(w io.Writer, s Summary) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}
