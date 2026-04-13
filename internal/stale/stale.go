// Package stale identifies drift results that have not changed between
// successive runs, flagging them as stale after a configurable age threshold.
package stale

import (
	"time"

	"github.com/driftwatch/internal/drift"
)

// Options controls staleness detection behaviour.
type Options struct {
	// After is the duration after which an unchanged drift result is considered stale.
	After time.Duration
	// Reference is the timestamp of the previous run. Results older than
	// Reference+After are marked stale.
	Reference time.Time
}

// Result wraps a drift.Result with staleness metadata.
type Result struct {
	drift.Result
	Stale     bool
	FirstSeen time.Time
	Age       time.Duration
}

// Mark annotates each result with staleness information. firstSeen maps a
// "service:field" key to the time the drift was first observed.
func Mark(results []drift.Result, firstSeen map[string]time.Time, opts Options) []Result {
	now := opts.Reference
	if now.IsZero() {
		now = time.Now().UTC()
	}

	out := make([]Result, 0, len(results))
	for _, r := range results {
		sr := Result{Result: r}
		if r.Drifted {
			key := staleKey(r)
			if t, ok := firstSeen[key]; ok {
				sr.FirstSeen = t
				sr.Age = now.Sub(t)
				if opts.After > 0 && sr.Age >= opts.After {
					sr.Stale = true
				}
			} else {
				sr.FirstSeen = now
			}
		}
		out = append(out, sr)
	}
	return out
}

// UpdateFirstSeen merges new drift results into an existing firstSeen map,
// recording the observation time for any key not already present.
func UpdateFirstSeen(existing map[string]time.Time, results []drift.Result, at time.Time) map[string]time.Time {
	if existing == nil {
		existing = make(map[string]time.Time)
	}
	for _, r := range results {
		if !r.Drifted {
			continue
		}
		key := staleKey(r)
		if _, ok := existing[key]; !ok {
			existing[key] = at
		}
	}
	return existing
}

func staleKey(r drift.Result) string {
	return r.Service + ":" + r.Field
}
