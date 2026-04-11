// Package filter provides post-detection filtering of drift results.
//
// It allows callers to narrow the output of [drift.Detect] by service name
// patterns and by drift presence, making it easy to surface only actionable
// findings in large deployments.
//
// Typical usage:
//
//	results, _ := drift.Detect(declared, live)
//	filtered := filter.Apply(results, filter.Options{
//		Services:    []string{"api", "worker"},
//		OnlyDrifted: true,
//	})
//
// Patterns are matched as case-insensitive substrings against the service
// name, so "api" matches both "api-gateway" and "internal-api".
package filter
