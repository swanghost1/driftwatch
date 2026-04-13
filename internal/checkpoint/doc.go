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
package checkpoint
