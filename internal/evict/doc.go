// Package evict provides age-based and status-based eviction of drift
// results.
//
// Eviction is a lightweight pre-processing step that discards stale or
// irrelevant results before they reach reporting, alerting, or storage
// stages. Two eviction modes are supported:
//
//   - Age eviction: results whose DetectedAt timestamp is older than
//     MaxAge are dropped. This prevents historical noise from polluting
//     live dashboards.
//
//   - Status eviction: when OnlyDrifted is set, clean results are
//     removed immediately, useful when downstream consumers only care
//     about active drift.
//
// Both modes may be combined. The Apply function is stateless and safe
// for concurrent use.
package evict
