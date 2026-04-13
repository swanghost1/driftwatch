// Package rollup provides aggregation of drift detection results into
// group-level summaries.
//
// Groups are defined by a caller-supplied key function that maps a service
// name to an arbitrary string (e.g. a Kubernetes namespace, an environment
// prefix, or a tag value). Results that map to an empty string are collected
// under the synthetic group name "(untagged)".
//
// Typical usage:
//
//	groups := rollup.ByGroup(results, func(svc string) string {
//		// derive group from service name or metadata
//		return namespaceOf(svc)
//	})
//	rollup.Write(os.Stdout, groups)
package rollup
