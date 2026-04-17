// Package bucketize groups drift detection results into discrete time-based
// buckets (hourly, daily, or weekly) to support trend analysis and
// visualisation of drift frequency over time.
//
// Usage:
//
//	buckets := bucketize.Apply(results, bucketize.Daily)
//	bucketize.Write(os.Stdout, buckets)
//
// Granularities:
//
//	Hourly  — one bucket per clock hour
//	Daily   — one bucket per calendar day (default)
//	Weekly  — one bucket per ISO week
package bucketize
