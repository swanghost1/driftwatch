// Package replay allows re-running drift detection against a previously
// saved snapshot, producing a fresh set of results without contacting live
// infrastructure.
package replay

import (
	"fmt"
	"time"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/snapshot"
)

// Options controls how a replay is executed.
type Options struct {
	// SnapshotPath is the file from which the snapshot store is loaded.
	SnapshotPath string
	// AsOf, when non-zero, skips snapshots captured after this time.
	AsOf time.Time
}

// Result holds the outcome of a single replayed service.
type Result struct {
	drift.Result
	// ReplayedAt is when the replay was performed.
	ReplayedAt time.Time
}

// Run loads a snapshot store and re-runs drift detection against every
// entry, returning one Result per service.
func Run(declared []snapshot.Entry, opts Options) ([]Result, error) {
	store, err := snapshot.Load(opts.SnapshotPath)
	if err != nil {
		return nil, fmt.Errorf("replay: load snapshot: %w", err)
	}

	now := time.Now().UTC()
	var results []Result

	for _, decl := range declared {
		live, ok := store.FindByName(decl.Name)
		if !ok {
			continue
		}
		if !opts.AsOf.IsZero() && live.CapturedAt.After(opts.AsOf) {
			continue
		}
		dr := drift.Detect(decl, live)
		results = append(results, Result{
			Result:     dr,
			ReplayedAt: now,
		})
	}
	return results, nil
}
