// Package aggregate provides field-level aggregation of drift results
// across all monitored services.
//
// Rather than viewing drift per service, aggregate groups results by
// the field that drifted (e.g. "image", "replicas", "env.PORT") and
// computes totals and drift rates across the entire service fleet.
//
// Typical usage:
//
//	summaries := aggregate.ByField(results)
//	_ = aggregate.Write(os.Stdout, summaries)
//
// This is useful for identifying which configuration fields are most
// prone to drift across a deployment, helping teams prioritise
// remediation or policy enforcement efforts.
package aggregate
