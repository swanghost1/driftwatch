package watch_test

import (
	"context"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/driftwatch/internal/watch"
)

func writeTempConfig(t *testing.T, dir, content string) string {
	t.Helper()
	p := filepath.Join(dir, "driftwatch.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempConfig: %v", err)
	}
	return p
}

func TestWatch_CallsHandlerOnWrite(t *testing.T) {
	dir := t.TempDir()
	path := writeTempConfig(t, dir, "version: 1\n")

	var calls atomic.Int32
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	opts := watch.Options{Debounce: 50 * time.Millisecond}

	errCh := make(chan error, 1)
	go func() {
		errCh <- watch.Watch(ctx, path, opts, func(p string) error {
			calls.Add(1)
			return nil
		})
	}()

	time.Sleep(100 * time.Millisecond)
	if err := os.WriteFile(path, []byte("version: 2\n"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if calls.Load() >= 1 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	cancel()
	if err := <-errCh; err != nil {
		t.Fatalf("Watch returned error: %v", err)
	}
	if calls.Load() < 1 {
		t.Errorf("expected handler to be called at least once, got %d", calls.Load())
	}
}

func TestWatch_CancelStopsWatcher(t *testing.T) {
	dir := t.TempDir()
	path := writeTempConfig(t, dir, "version: 1\n")

	ctx, cancel := context.WithCancel(context.Background())
	opts := watch.Options{Debounce: 50 * time.Millisecond}

	errCh := make(chan error, 1)
	go func() {
		errCh <- watch.Watch(ctx, path, opts, func(string) error { return nil })
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("expected nil error on cancel, got: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Watch did not stop after context cancellation")
	}
}

func TestWatch_InvalidPath_ReturnsError(t *testing.T) {
	ctx := context.Background()
	err := watch.Watch(ctx, "/nonexistent/path/driftwatch.yaml", watch.Options{}, func(string) error { return nil })
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}
