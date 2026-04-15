// Package quota enforces per-service drift field limits for driftwatch runs.
//
// A quota defines the maximum number of drifted fields that a single service
// is allowed to report before a violation is raised. This is useful for
// detecting runaway services that are heavily out of sync with their declared
// state, and for alerting operators to services that may require immediate
// attention.
//
// Usage:
//
//	opts := quota.Options{MaxDriftedFields: 3}
//	result := quota.Evaluate(driftResults, opts)
//	if result.Exceeded() {
//		quota.Write(os.Stdout, result)
//	}
//
// A MaxDriftedFields value of zero disables enforcement entirely.
package quota
