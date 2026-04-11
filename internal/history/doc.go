// Package history records the results of each drift detection run to disk
// and provides utilities for listing past runs and analysing drift trends
// over time.
//
// # Storage layout
//
// Each run is persisted as a JSON file named by its UTC timestamp:
//
//	<history-dir>/20060102T150405Z.json
//
// Files are written atomically via os.Create and are sorted lexicographically
// by filename, which preserves chronological order.
//
// # Usage
//
//	store := history.NewStore(".driftwatch/history")
//
//	// Record results after a detection run.
//	if err := store.Record(results); err != nil { ... }
//
//	// List all past runs and render a trend table.
//	entries, err := store.List()
//	summaries := history.Summarise(entries)
//	history.WriteTrend(os.Stdout, summaries)
package history
