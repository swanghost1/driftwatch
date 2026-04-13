// Package stale tracks how long individual drift results have persisted
// across successive driftwatch runs and marks them as stale once they
// exceed a configurable age threshold.
//
// Typical usage:
//
//	store := stale.NewStore(".driftwatch/stale.json")
//	firstSeen, _ := store.Load()
//
//	// After running drift detection:
//	firstSeen = stale.UpdateFirstSeen(firstSeen, results, time.Now().UTC())
//	_ = store.Save(firstSeen)
//
//	marked := stale.Mark(results, firstSeen, stale.Options{
//		After:     48 * time.Hour,
//		Reference: time.Now().UTC(),
//	})
//
// A Result is considered stale when its drift has been continuously present
// for longer than Options.After. Non-drifted results are never marked stale.
package stale
