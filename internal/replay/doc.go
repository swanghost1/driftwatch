// Package replay re-runs drift detection against a previously saved snapshot
// store, allowing operators to audit historical state without needing access
// to live infrastructure.
//
// Typical usage:
//
//	results, err := replay.Run(declaredEntries, replay.Options{
//		SnapshotPath: ".driftwatch/snapshots/latest.json",
//	})
//	if err != nil {
//		log.Fataln//	replay.WriteText(os.Stdout, results)
//
// The AsOf option can be used to restrict replay to snapshots captured
// before a given point in time, enabling point-in-time comparisons.
package replay
