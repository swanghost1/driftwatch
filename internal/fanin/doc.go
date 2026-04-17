// Package fanin consolidates results from multiple drift detection sources
// into a single ordered slice.
//
// When driftwatch runs detection against several config sources in parallel,
// each pass produces an independent []drift.Result. fanin.Apply merges those
// slices, optionally deduplicating by (service, field) key and sorting so
// that drifted entries appear first — making downstream reporting consistent
// regardless of source ordering.
//
// Example:
//
//	results := fanin.Apply(
//		[][]drift.Result{passA, passB},
//		fanin.DefaultOptions(),
//	)
package fanin
