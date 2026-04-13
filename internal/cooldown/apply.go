package cooldown

import (
	"time"

	"github.com/example/driftwatch/internal/drift"
)

// Options controls which results are suppressed by the cooldown store.
type Options struct {
	// Period is the quiet window during which repeated drift for the same
	// service is suppressed.
	Period time.Duration

	// RecordAll, when true, calls Record for every drifted service regardless
	// of whether it was suppressed.
	RecordAll bool
}

// Apply filters results so that services still within their cooldown period
// are removed from the returned slice. Drifted services that pass (or are
// newly seen) are recorded in the store.
func Apply(store *Store, results []drift.Result, opts Options) ([]drift.Result, error) {
	var out []drift.Result
	for _, r := range results {
		if !r.Drifted {
			out = append(out, r)
			continue
		}
		if store.IsCoolingDown(r.Service, opts.Period) {
			if opts.RecordAll {
				if err := store.Record(r.Service); err != nil {
					return nil, err
				}
			}
			continue
		}
		if err := store.Record(r.Service); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}
