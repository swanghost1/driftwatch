// Package jitter provides utilities for adding randomised jitter to durations.
//
// Jitter is commonly applied to scheduled intervals to prevent thundering-herd
// problems when many driftwatch instances run concurrently against the same
// infrastructure. By spreading check times slightly, load on downstream APIs
// is smoothed out.
//
// Example:
//
//	opts := jitter.DefaultOptions() // Factor: 0.1
//	interval := jitter.Apply(30*time.Second, opts)
//	// interval is somewhere in [30s, 33s]
package jitter
