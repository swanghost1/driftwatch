// Package headroom computes how much "drift budget" remains before a
// configured threshold is breached, giving operators an early-warning
// indicator before policy limits are hit.
package headroom

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/driftwatch/internal/drift"
)

// Options controls how headroom is calculated.
type Options struct {
	// MaxDrifted is the absolute upper limit on drifted services.
	// Zero means unlimited (headroom is always 100 %).
	MaxDrifted int

	// MaxRatePct is the upper limit expressed as a percentage (0-100).
	// Zero means unlimited.
	MaxRatePct float64
}

// DefaultOptions returns conservative defaults.
func DefaultOptions() Options {
	return Options{
		MaxDrifted: 0,
		MaxRatePct: 0,
	}
}

// Report describes the remaining headroom.
type Report struct {
	Total       int     `json:"total"`
	Drifted     int     `json:"drifted"`
	Clean       int     `json:"clean"`
	DriftRatePct float64 `json:"drift_rate_pct"`

	// CountHeadroom is MaxDrifted - Drifted; -1 when no limit is set.
	CountHeadroom int `json:"count_headroom"`

	// RateHeadroom is MaxRatePct - DriftRatePct; -1 when no limit is set.
	RateHeadroom float64 `json:"rate_headroom_pct"`
}

// Compute derives a headroom Report from a slice of drift results.
func Compute(results []drift.Result, opts Options) Report {
	total := len(results)
	drifted := 0
	for _, r := range results {
		if r.Drifted {
			drifted++
		}
	}

	var rate float64
	if total > 0 {
		rate = float64(drifted) / float64(total) * 100
	}

	countHeadroom := -1
	if opts.MaxDrifted > 0 {
		countHeadroom = opts.MaxDrifted - drifted
	}

	rateHeadroom := -1.0
	if opts.MaxRatePct > 0 {
		rateHeadroom = opts.MaxRatePct - rate
	}

	return Report{
		Total:         total,
		Drifted:       drifted,
		Clean:         total - drifted,
		DriftRatePct:  rate,
		CountHeadroom: countHeadroom,
		RateHeadroom:  rateHeadroom,
	}
}

// WriteText writes a human-readable headroom summary to w.
func WriteText(w io.Writer, r Report) {
	fmt.Fprintf(w, "Headroom Report\n")
	fmt.Fprintf(w, "  Total services : %d\n", r.Total)
	fmt.Fprintf(w, "  Drifted        : %d\n", r.Drifted)
	fmt.Fprintf(w, "  Clean          : %d\n", r.Clean)
	fmt.Fprintf(w, "  Drift rate     : %.1f%%\n", r.DriftRatePct)
	if r.CountHeadroom >= 0 {
		fmt.Fprintf(w, "  Count headroom : %d\n", r.CountHeadroom)
	}
	if r.RateHeadroom >= 0 {
		fmt.Fprintf(w, "  Rate headroom  : %.1f%%\n", r.RateHeadroom)
	}
}

// WriteJSON writes the Report as a JSON object to w.
func WriteJSON(w io.Writer, r Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
