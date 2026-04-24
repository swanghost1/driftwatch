// Package evict removes drift results that exceed a maximum age or
// fall outside a rolling retention window, keeping result sets lean
// for downstream consumers.
package evict

import (
	"fmt"
	"io"
	"time"

	"github.com/example/driftwatch/internal/drift"
)

// DefaultOptions returns an Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		MaxAge:      72 * time.Hour,
		OnlyDrifted: false,
	}
}

// Options controls which results are evicted.
type Options struct {
	// MaxAge is the maximum age of a result's DetectedAt timestamp.
	// Results older than this are removed. Zero disables age eviction.
	MaxAge time.Duration

	// OnlyDrifted, when true, evicts clean (non-drifted) results regardless
	// of age, retaining only results where Drifted == true.
	OnlyDrifted bool

	// Reference is the time used as "now" when computing age.
	// Defaults to time.Now() when zero.
	Reference time.Time
}

// Apply filters results according to opts, returning the survivors.
func Apply(results []drift.Result, opts Options) []drift.Result {
	if len(results) == 0 {
		return nil
	}

	now := opts.Reference
	if now.IsZero() {
		now = time.Now()
	}

	out := results[:0:0]
	for _, r := range results {
		if opts.OnlyDrifted && !r.Drifted {
			continue
		}
		if opts.MaxAge > 0 && !r.DetectedAt.IsZero() {
			if now.Sub(r.DetectedAt) > opts.MaxAge {
				continue
			}
		}
		out = append(out, r)
	}
	return out
}

// Write prints a short eviction summary to w.
func Write(w io.Writer, before, after int) {
	evicted := before - after
	fmt.Fprintf(w, "evict: %d result(s) removed, %d retained\n", evicted, after)
}
