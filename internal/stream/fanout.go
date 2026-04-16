package stream

import (
	"sync"

	"github.com/example/driftwatch/internal/drift"
)

// Fanout distributes results from a single source channel to multiple
// subscriber channels. Each subscriber receives every result.
type Fanout struct {
	mu   sync.Mutex
	subs []chan drift.Result
	buf  int
}

// NewFanout creates a Fanout where each subscriber channel has the given
// buffer size.
func NewFanout(bufSize int) *Fanout {
	return &Fanout{buf: bufSize}
}

// Subscribe returns a new channel that will receive all future results.
func (f *Fanout) Subscribe() <-chan drift.Result {
	ch := make(chan drift.Result, f.buf)
	f.mu.Lock()
	f.subs = append(f.subs, ch)
	f.mu.Unlock()
	return ch
}

// Publish sends r to all current subscribers. Slow subscribers that have a
// full buffer are skipped to avoid blocking the publisher.
func (f *Fanout) Publish(r drift.Result) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, ch := range f.subs {
		select {
		case ch <- r:
		default:
			// subscriber is not keeping up; skip this result for them
		}
	}
}

// Close closes all subscriber channels signalling end of stream.
func (f *Fanout) Close() {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, ch := range f.subs {
		close(ch)
	}
	f.subs = nil
}

// Run reads from src and fans out to all subscribers until src is closed.
func (f *Fanout) Run(src <-chan drift.Result) {
	for r := range src {
		f.Publish(r)
	}
	f.Close()
}
