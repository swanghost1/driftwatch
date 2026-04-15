// Package cap limits the maximum number of drift results returned per service,
// preventing noisy services from dominating reports.
package cap

import (
	"fmt"
	"io"

	"github.com/driftwatch/internal/drift"
)

// Options configures the cap behaviour.
type Options struct {
	// MaxPerService is the maximum number of drift results kept per service.
	// Zero means no cap is applied.
	MaxPerService int
}

// DefaultOptions returns sensible defaults (no cap).
func DefaultOptions() Options {
	return Options{MaxPerService: 0}
}

// Apply returns a copy of results with at most MaxPerService entries per
// service name. Results within the cap are kept in their original order.
// If opts.MaxPerService is zero or negative, results are returned unchanged.
func Apply(results []drift.Result, opts Options) []drift.Result {
	if opts.MaxPerService <= 0 {
		return results
	}

	counts := make(map[string]int)
	out := make([]drift.Result, 0, len(results))

	for _, r := range results {
		if counts[r.Service] < opts.MaxPerService {
			out = append(out, r)
			counts[r.Service]++
		}
	}

	return out
}

// Write prints a human-readable summary of the cap operation to w.
func Write(w io.Writer, before, after []drift.Result, opts Options) {
	dropped := len(before) - len(after)
	fmt.Fprintf(w, "cap: max_per_service=%d  before=%d  after=%d  dropped=%d\n",
		opts.MaxPerService, len(before), len(after), dropped)
}
