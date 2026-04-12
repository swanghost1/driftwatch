// Package watch provides file-system change detection for driftwatch
// configuration files.
//
// It wraps fsnotify with debouncing so that rapid successive writes
// (e.g. from editors that save atomically) result in a single handler
// invocation rather than a burst of calls.
//
// Typical usage:
//
//	opts := watch.Options{Debounce: 500 * time.Millisecond}
//	err := watch.Watch(ctx, "driftwatch.yaml", opts, func(path string) error {
//	    results, err := runner.Run(cfg)
//	    // handle results ...
//	    return err
//	})
//
// Watch blocks until the provided context is cancelled or a fatal
// watcher error occurs. A nil error is returned on clean shutdown.
package watch
