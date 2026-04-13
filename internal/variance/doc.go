// Package variance analyses per-service drift counts across historical runs
// to compute statistical variance metrics.
//
// # Overview
//
// For each service observed across multiple runs, variance computes the mean
// and standard deviation of the number of drifted fields. A service is flagged
// as anomalous when its latest drift count exceeds mean + 2*stddev, indicating
// a statistically significant increase relative to its own history.
//
// # Usage
//
//	store := variance.NewStore(".driftwatch/variance")
//	if err := store.Record(results); err != nil { ... }
//
//	entries, _ := store.Load(20)
//	history := make([][]drift.Result, len(entries))
//	for i, e := range entries { history[i] = e.Results }
//
//	variances := variance.Compute(history)
//	variance.Write(os.Stdout, variances)
//
// At least three historical samples are required before a service can be
// marked anomalous, avoiding false positives on sparse data.
package variance
