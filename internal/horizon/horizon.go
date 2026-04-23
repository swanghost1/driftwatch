// Package horizon provides a look-ahead window that predicts whether drift
// is likely to exceed a threshold within a future time horizon, based on
// the recent rate of change observed in historical run data.
package horizon

import (
	"fmt"
	"io"
	"time"

	"github.com/example/driftwatch/internal/history"
)

// DefaultOptions returns sensible defaults for horizon prediction.
func DefaultOptions() Options {
	return Options{
		Window:    7 * 24 * time.Hour,
		Horizon:   24 * time.Hour,
		Threshold: 0.5,
	}
}

// Options configures the horizon prediction.
type Options struct {
	// Window is how far back in history to sample drift rate from.
	Window time.Duration
	// Horizon is how far into the future to project.
	Horizon time.Duration
	// Threshold is the drift rate (0–1) that, if predicted to be reached,
	// triggers a warning.
	Threshold float64
}

// Prediction is the result of a horizon evaluation.
type Prediction struct {
	CurrentRate   float64
	ProjectedRate float64
	Horizon       time.Duration
	WillExceed    bool
}

// Evaluate projects the drift rate over the configured horizon using entries
// from the provided history store.
func Evaluate(entries []history.Entry, now time.Time, opts Options) *Prediction {
	if len(entries) == 0 {
		return nil
	}

	cutoff := now.Add(-opts.Window)
	var recent []history.Entry
	for _, e := range entries {
		if e.RecordedAt.After(cutoff) {
			recent = append(recent, e)
		}
	}
	if len(recent) == 0 {
		return nil
	}

	var totalRate float64
	for _, e := range recent {
		if e.Total > 0 {
			totalRate += float64(e.Drifted) / float64(e.Total)
		}
	}
	currentRate := totalRate / float64(len(recent))

	// Simple linear projection: scale rate by horizon/window ratio.
	ratio := opts.Horizon.Seconds() / opts.Window.Seconds()
	projected := currentRate + (currentRate * ratio)
	if projected > 1.0 {
		projected = 1.0
	}

	return &Prediction{
		CurrentRate:   currentRate,
		ProjectedRate: projected,
		Horizon:       opts.Horizon,
		WillExceed:    projected >= opts.Threshold,
	}
}

// Write renders the prediction as human-readable text to w.
func Write(w io.Writer, p *Prediction) {
	if p == nil {
		fmt.Fprintln(w, "horizon: no data available")
		return
	}
	status := "OK"
	if p.WillExceed {
		status = "WARNING"
	}
	fmt.Fprintf(w, "horizon (%s): current_rate=%.1f%% projected_rate=%.1f%% status=%s\n",
		p.Horizon.String(),
		p.CurrentRate*100,
		p.ProjectedRate*100,
		status,
	)
}
