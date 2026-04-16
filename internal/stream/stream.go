// Package stream provides live streaming of drift results as they are
// produced, writing each result to an io.Writer incrementally rather than
// buffering the full slice.
package stream

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/example/driftwatch/internal/drift"
)

// Options controls streaming behaviour.
type Options struct {
	// Format is either "text" or "json".
	Format string
	// FlushInterval is how often a heartbeat line is written when there is no
	// activity. Zero disables heartbeats.
	FlushInterval time.Duration
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Format:        "text",
		FlushInterval: 5 * time.Second,
	}
}

// Write streams each drift.Result from ch to w as it arrives.
// It returns when ch is closed or ctx is done.
func Write(w io.Writer, ch <-chan drift.Result, opts Options) error {
	if opts.Format == "json" {
		return writeJSON(w, ch, opts)
	}
	return writeText(w, ch, opts)
}

func writeText(w io.Writer, ch <-chan drift.Result, opts Options) error {
	var ticker *time.Ticker
	var tick <-chan time.Time
	if opts.FlushInterval > 0 {
		ticker = time.NewTicker(opts.FlushInterval)
		defer ticker.Stop()
		tick = ticker.C
	}
	for {
		select {
		case r, ok := <-ch:
			if !ok {
				return nil
			}
			status := "OK"
			if r.Drifted {
				status = "DRIFT"
			}
			if _, err := fmt.Fprintf(w, "[%s] %-30s %-15s want=%s got=%s\n",
				status, r.Service, r.Field, r.Want, r.Got); err != nil {
				return err
			}
		case <-tick:
			if _, err := fmt.Fprintln(w, "# heartbeat"); err != nil {
				return err
			}
		}
	}
}

func writeJSON(w io.Writer, ch <-chan drift.Result, opts Options) error {
	enc := json.NewEncoder(w)
	var ticker *time.Ticker
	var tick <-chan time.Time
	if opts.FlushInterval > 0 {
		ticker = time.NewTicker(opts.FlushInterval)
		defer ticker.Stop()
		tick = ticker.C
	}
	for {
		select {
		case r, ok := <-ch:
			if !ok {
				return nil
			}
			if err := enc.Encode(r); err != nil {
				return err
			}
		case <-tick:
			if _, err := fmt.Fprintln(w, `{"heartbeat":true}`); err != nil {
				return err
			}
		}
	}
}
