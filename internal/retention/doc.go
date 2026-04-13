// Package retention implements pruning of file-based records produced by
// driftwatch subsystems such as history and audit.
//
// A Policy controls two independent limits:
//
//   - MaxAge: entries older than this duration are removed regardless of count.
//   - MaxEntries: if more than this many entries remain after age pruning, the
//     oldest are removed until the count is within the limit. A value of 0
//     disables the limit.
//
// Prune operates on a directory of JSON files and returns the list of paths
// that were deleted, making it easy to surface results to the operator via
// Write.
package retention
