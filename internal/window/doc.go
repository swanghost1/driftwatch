// Package window provides time-based filtering for drift results.
//
// It allows callers to restrict a result set to those detected within a
// specific time range, defined by an optional lower bound (Since) and an
// optional upper bound (Until).
//
// Example usage:
//
//	filtered := window.Apply(results, window.Options{
//		Since: time.Now().Add(-24 * time.Hour),
//	})
//
// The Last helper is a convenience wrapper for rolling windows:
//
//	recent := window.Last(results, 6*time.Hour)
package window
