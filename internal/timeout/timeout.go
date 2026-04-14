package timeout

import (
	"errors"
	"time"

	"github.com/driftwatch/internal/drift"
)

// ErrTimeout is returned when a detection run exceeds its deadline.
var ErrTimeout = errors.New("timeout: detection exceeded deadline")

// Options controls timeout behaviour.
type Options struct {
	// Deadline is the maximum duration allowed for a detection run.
	// Zero means no timeout is applied.
	Deadline time.Duration

	// GracePeriod is added on top of Deadline before hard-cancellation.
	GracePeriod time.Duration
}

// DefaultOptions returns sensible defaults: 30 s deadline, 5 s grace.
func DefaultOptions() Options {
	return Options{
		Deadline:    30 * time.Second,
		GracePeriod: 5 * time.Second,
	}
}

// RunFunc is the signature of a detection function that timeout wraps.
type RunFunc func() ([]drift.Result, error)

// Apply runs fn and returns its results. If fn does not complete within
// opts.Deadline + opts.GracePeriod, Apply returns ErrTimeout.
// When opts.Deadline is zero the function is called directly with no
// timeout applied.
func Apply(opts Options, fn RunFunc) ([]drift.Result, error) {
	if opts.Deadline == 0 {
		return fn()
	}

	total := opts.Deadline + opts.GracePeriod

	type result struct {
		res []drift.Result
		err error
	}

	ch := make(chan result, 1)
	go func() {
		r, err := fn()
		ch <- result{r, err}
	}()

	select {
	case r := <-ch:
		return r.res, r.err
	case <-time.After(total):
		return nil, ErrTimeout
	}
}
