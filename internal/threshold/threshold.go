package threshold

import (
	"fmt"
	"io"

	"github.com/driftwatch/internal/drift"
)

// Options controls threshold evaluation behaviour.
type Options struct {
	// MinDriftCount triggers a breach when drifted results meet or exceed this value.
	MinDriftCount int
	// MinDriftRate triggers a breach when the drift rate (0–1) meets or exceeds this value.
	MinDriftRate float64
}

// DefaultOptions returns sensible defaults (no thresholds enforced).
func DefaultOptions() Options {
	return Options{}
}

// Result holds the outcome of a threshold evaluation.
type Result struct {
	Total      int
	Drifted    int
	DriftRate  float64
	Breached   bool
	Reason     string
}

// Evaluate checks whether the supplied drift results breach the configured thresholds.
func Evaluate(results []drift.Result, opts Options) Result {
	total := len(results)
	drifted := 0
	for _, r := range results {
		if r.Drifted {
			drifted++
		}
	}

	rate := 0.0
	if total > 0 {
		rate = float64(drifted) / float64(total)
	}

	r := Result{
		Total:     total,
		Drifted:   drifted,
		DriftRate: rate,
	}

	if opts.MinDriftCount > 0 && drifted >= opts.MinDriftCount {
		r.Breached = true
		r.Reason = fmt.Sprintf("drifted count %d meets threshold %d", drifted, opts.MinDriftCount)
		return r
	}

	if opts.MinDriftRate > 0 && rate >= opts.MinDriftRate {
		r.Breached = true
		r.Reason = fmt.Sprintf("drift rate %.2f meets threshold %.2f", rate, opts.MinDriftRate)
		return r
	}

	return r
}

// Write prints a human-readable threshold summary to w.
func Write(w io.Writer, r Result) {
	status := "OK"
	if r.Breached {
		status = "BREACHED"
	}
	fmt.Fprintf(w, "Threshold: %s\n", status)
	fmt.Fprintf(w, "  Total:     %d\n", r.Total)
	fmt.Fprintf(w, "  Drifted:   %d\n", r.Drifted)
	fmt.Fprintf(w, "  DriftRate: %.2f\n", r.DriftRate)
	if r.Breached {
		fmt.Fprintf(w, "  Reason:    %s\n", r.Reason)
	}
}
