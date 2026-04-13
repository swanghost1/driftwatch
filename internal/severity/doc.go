// Package severity classifies drift detection results into severity levels
// (none, low, medium, high, critical) based on which fields have drifted.
//
// Each field name can be assigned a severity level via a Config map.
// When multiple fields drift in a single service, the highest level wins.
//
// Typical usage:
//
//	cfg := severity.DefaultConfig()
//	classified := severity.Classify(driftResults, cfg)
//	_ = severity.Write(os.Stdout, classified)
//
// DefaultConfig provides sensible defaults:
//
//	image    → critical
//	replicas → high
//	env      → medium
//	port     → low
package severity
