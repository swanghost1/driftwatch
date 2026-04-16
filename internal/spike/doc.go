// Package spike provides drift-spike detection for driftwatch.
//
// A spike is defined as a run where the number of drifted services
// significantly exceeds the rolling mean of recent runs. This helps
// distinguish routine, low-level drift from sudden, unexpected
// configuration changes across many services.
//
// Basic usage:
//
//	store := spike.NewStore(".driftwatch/spike.json")
//	history, _ := store.Counts()
//	result := spike.Detect(currentDriftCount, history, spike.DefaultOptions())
//	if result.IsSpike {
//		spike.Write(os.Stdout, result)
//	}
//	_ = store.Record(currentDriftCount, time.Now())
package spike
