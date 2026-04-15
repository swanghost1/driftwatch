// Package checkpoint records and retrieves the last known scan outcome for
// each named driftwatch run.
//
// A checkpoint captures when the scan ran, how many services were evaluated,
// and how many were found to be drifted or clean. This allows operators and
// automation to quickly determine whether a re-scan is necessary without
// re-running the full detection pipeline.
//
// Checkpoints are stored as JSON files in a configurable directory, keyed by
// the scan name. Saving a checkpoint is idempotent — repeated saves overwrite
// the previous entry for that name.
//
// # File Layout
//
// Each checkpoint file is named "<scan-name>.json" and lives under the
// configured storage directory. The JSON structure mirrors the Checkpoint
// struct, including a UTC timestamp, service counts, and a drift summary.
//
// # Concurrency
//
// The package does not provide its own locking. Callers that write checkpoints
// from multiple goroutines for the same scan name should coordinate access
// externally to avoid partial writes.
package checkpoint
