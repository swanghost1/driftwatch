// Package variance tracks how much drift metrics deviate from a rolling
// baseline, helping surface services whose drift behaviour is changing.
package variance

import (
	"fmt"
	"io"
	"math"
	"sort"

	"github.com/driftwatch/internal/drift"
)

// ServiceVariance holds the computed variance statistics for a single service.
type ServiceVariance struct {
	Service  string
	Mean     float64
	StdDev   float64
	Samples  int
	Anomalous bool // true when latest drift count exceeds mean + 2*stddev
}

// Compute calculates per-service drift variance from a slice of historical
// run results. Each element of history represents one run's results.
func Compute(history [][]drift.Result) []ServiceVariance {
	if len(history) == 0 {
		return nil
	}

	// Accumulate drift counts per service across runs.
	counts := map[string][]float64{}
	for _, run := range history {
		perService := map[string]int{}
		for _, r := range run {
			if r.Drifted {
				perService[r.Service]++
			} else {
				// ensure service appears even with zero drift
				if _, ok := perService[r.Service]; !ok {
					perService[r.Service] = 0
				}
			}
		}
		for svc, n := range perService {
			counts[svc] = append(counts[svc], float64(n))
		}
	}

	svcs := make([]string, 0, len(counts))
	for svc := range counts {
		svcs = append(svcs, svc)
	}
	sort.Strings(svcs)

	result := make([]ServiceVariance, 0, len(svcs))
	for _, svc := range svcs {
		samples := counts[svc]
		mean, stddev := stats(samples)
		latest := samples[len(samples)-1]
		result = append(result, ServiceVariance{
			Service:   svc,
			Mean:      mean,
			StdDev:    stddev,
			Samples:   len(samples),
			Anomalous: len(samples) >= 3 && latest > mean+2*stddev,
		})
	}
	return result
}

// Write renders variance statistics in a human-readable table.
func Write(w io.Writer, variances []ServiceVariance) {
	fmt.Fprintf(w, "%-30s %8s %8s %8s %s\n", "SERVICE", "MEAN", "STDDEV", "SAMPLES", "ANOMALOUS")
	fmt.Fprintf(w, "%s\n", fmt.Sprintf("%s", repeatChar('-', 65)))
	for _, v := range variances {
		anom := ""
		if v.Anomalous {
			anom = "YES"
		}
		fmt.Fprintf(w, "%-30s %8.2f %8.2f %8d %s\n",
			v.Service, v.Mean, v.StdDev, v.Samples, anom)
	}
}

func stats(vals []float64) (mean, stddev float64) {
	if len(vals) == 0 {
		return 0, 0
	}
	for _, v := range vals {
		mean += v
	}
	mean /= float64(len(vals))
	for _, v := range vals {
		diff := v - mean
		stddev += diff * diff
	}
	stddev = math.Sqrt(stddev / float64(len(vals)))
	return mean, stddev
}

func repeatChar(c rune, n int) string {
	out := make([]rune, n)
	for i := range out {
		out[i] = c
	}
	return string(out)
}
