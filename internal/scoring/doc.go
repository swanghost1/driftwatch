// Package scoring provides drift health scoring for driftwatch.
//
// A score is a value in the range [0, 100] that represents the proportion of
// services that are clean (not drifted). A score of 100 indicates that every
// service matches its declared state; a score of 0 means every service has
// drifted.
//
// Scores are also mapped to letter grades:
//
//	95–100  A  Excellent
//	80–94   B  Good
//	65–79   C  Fair
//	50–64   D  Poor
//	0–49    F  Critical
//
// The Trend type allows callers to accumulate scores over time and determine
// whether the overall drift health is improving, degrading, or stable.
package scoring
