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
// # Computing a Score
//
// Use [Compute] to derive a score from a snapshot of service states:
//
//	score := scoring.Compute(total, drifted)
//	grade := score.Grade()
//
// # Tracking Trends
//
// The [Trend] type allows callers to accumulate scores over time and determine
// whether the overall drift health is improving, degrading, or stable:
//
//	var t scoring.Trend
//	t.Record(score)
//	fmt.Println(t.Direction()) // Improving, Degrading, or Stable
package scoring
