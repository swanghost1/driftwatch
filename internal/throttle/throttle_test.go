package throttle_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/driftwatch/internal/throttle"
)

func storeAt(t *testing.T, name string) *throttle.Store {
	t.Helper()
	return throttle.NewStore(filepath.Join(t.TempDir(), name))
}

func TestCheck_NoStateFile_Passes(t *testing.T) {
	s := storeAt(t, "state.json")
	if err := s.Check(time.Minute); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRecord_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	s := throttle.NewStore(filepath.Join(dir, "state.json"))
	if err := s.Record(); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "state.json")); err != nil {
		t.Fatalf("state file not created: %v", err)
	}
}

func TestRecord_CreatesIntermediateDirectories(t *testing.T) {
	dir := t.TempDir()
	s := throttle.NewStore(filepath.Join(dir, "a", "b", "state.json"))
	if err := s.Record(); err != nil {
		t.Fatalf("Record: %v", err)
	}
}

func TestCheck_WithinCooldown_ReturnsThrottled(t *testing.T) {
	s := storeAt(t, "state.json")
	if err := s.Record(); err != nil {
		t.Fatalf("Record: %v", err)
	}
	err := s.Check(time.Hour)
	if !errors.Is(err, throttle.ErrThrottled) {
		t.Fatalf("expected ErrThrottled, got %v", err)
	}
}

func TestCheck_CooldownElapsed_Passes(t *testing.T) {
	s := storeAt(t, "state.json")
	if err := s.Record(); err != nil {
		t.Fatalf("Record: %v", err)
	}
	// A zero cooldown means the window is always elapsed.
	if err := s.Check(0); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCheck_AfterCooldown_Passes(t *testing.T) {
	s := storeAt(t, "state.json")
	if err := s.Record(); err != nil {
		t.Fatalf("Record: %v", err)
	}
	// Cooldown of 1ns is effectively already elapsed by the time we check.
	time.Sleep(2 * time.Millisecond)
	if err := s.Check(time.Nanosecond); err != nil {
		t.Fatalf("expected no error after cooldown, got %v", err)
	}
}
