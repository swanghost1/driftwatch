// Package spike detects sudden increases in drift count relative to a
// rolling baseline, flagging runs where drift has spiked unexpectedly.
package spike

import (
	"fmt"
	"io"
	"math"
)

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	return Options{
		WindowSize:    5,
		ThresholdMult: 2.0,
	}
}

// Options controls spike detection behaviour.
type Options struct {
	// WindowSize is the number of historical drift counts to average.
	WindowSize int
	// ThresholdMult is the multiplier applied to the rolling mean; a run
	// whose drift count exceeds mean*ThresholdMult is considered a spike.
	ThresholdMult float64
}

// Result describes whether the current run is a spike.
type Result struct {
	Current    int
	Mean       float64
	Threshold  float64
	IsSpike    bool
}

// Detect compares current drift count against the rolling mean of
// historical counts. history should be ordered oldest-first.
func Detect(current int, history []int, opts Options) Result {
	if opts.WindowSize <= 0 {
		opts.WindowSize = DefaultOptions().WindowSize
	}
	if opts.ThresholdMult <= 0 {
		opts.ThresholdMult = DefaultOptions().ThresholdMult
	}

	win := history
	if len(win) > opts.WindowSize {
		win = win[len(win)-opts.WindowSize:]
	}

	var mean float64
	if len(win) > 0 {
		sum := 0
		for _, v := range win {
			sum += v
		}
		mean = float64(sum) / float64(len(win))
	}

	threshold := math.Ceil(mean * opts.ThresholdMult)
	isSpike := len(win) > 0 && float64(current) >= threshold

	return Result{
		Current:   current,
		Mean:      mean,
		Threshold: threshold,
		IsSpike:   isSpike,
	}
}

// Write writes a human-readable spike report to w.
func Write(w io.Writer, r Result) {
	status := "no"
	if r.IsSpike {
		status = "YES"
	}
	fmt.Fprintf(w, "spike detected : %s\n", status)
	fmt.Fprintf(w, "current drift  : %d\n", r.Current)
	fmt.Fprintf(w, "rolling mean   : %.2f\n", r.Mean)
	fmt.Fprintf(w, "threshold      : %.0f\n", r.Threshold)
}
