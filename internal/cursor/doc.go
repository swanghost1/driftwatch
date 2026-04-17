// Package cursor tracks the last-processed position within a named result
// stream across incremental driftwatch runs.
//
// Each named cursor records:
//   - Offset    – number of results consumed so far
//   - LastSeen  – wall-clock time of the most recent advance
//   - RunCount  – how many times the cursor has been advanced
//
// Usage:
//
//	store := cursor.NewStore(".driftwatch/cursors")
//	st, err := store.Load("my-service")
//	if errors.Is(err, cursor.ErrNotFound) {
//		st = cursor.State{Name: "my-service"}
//	}
//	st = cursor.Advance(st, len(newResults))
//	store.Save(st)
package cursor
