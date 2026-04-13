package cooldown_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/cooldown"
)

func newStore(t *testing.T) *cooldown.Store {
	t.Helper()
	path := filepath.Join(t.TempDir(), "cooldown.json")
	s, err := cooldown.NewStore(path)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func TestIsCoolingDown_NoEntry_ReturnsFalse(t *testing.T) {
	s := newStore(t)
	if s.IsCoolingDown("svc-a", time.Minute) {
		t.Fatal("expected false for unknown service")
	}
}

func TestRecord_ThenIsCoolingDown_ReturnsTrue(t *testing.T) {
	s := newStore(t)
	if err := s.Record("svc-a"); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if !s.IsCoolingDown("svc-a", time.Hour) {
		t.Fatal("expected service to be cooling down")
	}
}

func TestIsCoolingDown_ExpiredPeriod_ReturnsFalse(t *testing.T) {
	s := newStore(t)
	if err := s.Record("svc-b"); err != nil {
		t.Fatalf("Record: %v", err)
	}
	// period of zero — always expired
	if s.IsCoolingDown("svc-b", 0) {
		t.Fatal("expected false for zero period")
	}
}

func TestRecord_CreatesIntermediateDirectories(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a", "b", "cooldown.json")
	s, err := cooldown.NewStore(path)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	if err := s.Record("svc-c"); err != nil {
		t.Fatalf("Record: %v", err)
	}
}

func TestNewStore_LoadsExistingEntries(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cooldown.json")
	s1, _ := cooldown.NewStore(path)
	_ = s1.Record("svc-d")

	s2, err := cooldown.NewStore(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if !s2.IsCoolingDown("svc-d", time.Hour) {
		t.Fatal("expected reloaded store to recognise svc-d")
	}
}

func TestIsCoolingDown_IndependentServices(t *testing.T) {
	s := newStore(t)
	_ = s.Record("svc-x")
	if s.IsCoolingDown("svc-y", time.Hour) {
		t.Fatal("svc-y should not be cooling down")
	}
}
