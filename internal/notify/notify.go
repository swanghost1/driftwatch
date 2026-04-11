// Package notify provides alerting support for drift detection results.
// It allows callers to send notifications when drift is detected, supporting
// configurable thresholds and multiple output channels.
package notify

import (
	"fmt"
	"io"
	"strings"

	"github.com/example/driftwatch/internal/drift"
)

// Channel represents a notification destination.
type Channel string

const (
	ChannelStdout Channel = "stdout"
	ChannelLog    Channel = "log"
)

// Options configures notification behaviour.
type Options struct {
	// Channel is the destination for notifications.
	Channel Channel
	// MinDriftCount is the minimum number of drifted services required to
	// trigger a notification. Zero means notify on any drift.
	MinDriftCount int
}

// Notifier sends drift alerts to a configured channel.
type Notifier struct {
	opts Options
	w    io.Writer
}

// New returns a Notifier that writes to w using opts.
func New(w io.Writer, opts Options) *Notifier {
	return &Notifier{opts: opts, w: w}
}

// Notify evaluates results and writes an alert if the drift threshold is met.
// It returns true if a notification was sent.
func (n *Notifier) Notify(results []drift.Result) (bool, error) {
	drifted := countDrifted(results)
	if drifted < n.opts.MinDriftCount || (n.opts.MinDriftCount == 0 && drifted == 0) {
		return false, nil
	}

	services := driftedNames(results)
	msg := fmt.Sprintf(
		"[driftwatch] DRIFT DETECTED: %d service(s) out of sync: %s\n",
		drifted,
		strings.Join(services, ", "),
	)

	if _, err := fmt.Fprint(n.w, msg); err != nil {
		return false, fmt.Errorf("notify: write failed: %w", err)
	}
	return true, nil
}

func countDrifted(results []drift.Result) int {
	count := 0
	for _, r := range results {
		if r.Drifted {
			count++
		}
	}
	return count
}

func driftedNames(results []drift.Result) []string {
	var names []string
	for _, r := range results {
		if r.Drifted {
			names = append(names, r.Service)
		}
	}
	return names
}
