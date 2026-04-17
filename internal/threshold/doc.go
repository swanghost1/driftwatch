// Package threshold evaluates drift results against configurable count and
// rate thresholds, producing a structured Result that callers can inspect or
// surface as a non-zero exit code.
//
// Usage:
//
//	opts := threshold.Options{
//	    MinDriftCount: 3,
//	    MinDriftRate:  0.5,
//	}
//	r := threshold.Evaluate(results, opts)
//	if r.Breached {
//	    threshold.Write(os.Stderr, r)
//	    os.Exit(1)
//	}
package threshold
