// Package watch provides file-system watching for driftwatch config files,
// triggering re-evaluation whenever the declared state changes on disk.
package watch

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Handler is called each time the watched config file changes.
type Handler func(path string) error

// Options controls watcher behaviour.
type Options struct {
	// Debounce is the minimum quiet period after the last event before
	// Handler is invoked. Defaults to 500 ms when zero.
	Debounce time.Duration
}

// Watch monitors path and calls h after each stable write event.
// It blocks until ctx is cancelled or a fatal error occurs.
//
// The handler is debounced: rapid successive file events (e.g. from editors
// that write via a temp-file rename) are collapsed into a single call.
func Watch(ctx context.Context, path string, opts Options, h Handler) error {
	if opts.Debounce == 0 {
		opts.Debounce = 500 * time.Millisecond
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("watch: create watcher: %w", err)
	}
	defer w.Close()

	if err := w.Add(path); err != nil {
		return fmt.Errorf("watch: add path %q: %w", path, err)
	}

	var debounce *time.Timer

	for {
		select {
		case <-ctx.Done():
			if debounce != nil {
				debounce.Stop()
			}
			return nil

		case event, ok := <-w.Events:
			if !ok {
				return nil
			}
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
				if debounce != nil {
					debounce.Stop()
				}
				debounce = time.AfterFunc(opts.Debounce, func() {
					if err := h(path); err != nil {
						log.Printf("watch: handler error: %v", err)
					}
				})
			}

		case err, ok := <-w.Errors:
			if !ok {
				return nil
			}
			log.Printf("watch: watcher error: %v", err)
		}
	}
}
