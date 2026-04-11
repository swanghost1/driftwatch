package schedule_test

import (
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/schedule"
)

func TestDuration_KnownIntervals(t *testing.T) {
	cases := []struct {
		interval schedule.Interval
		want     time.Duration
	}{
		{schedule.IntervalHourly, time.Hour},
		{schedule.IntervalDaily, 24 * time.Hour},
		{schedule.IntervalWeekly, 7 * 24 * time.Hour},
	}
	for _, tc := range cases {
		s := schedule.Schedule{Interval: tc.interval}
		got, err := s.Duration()
		if err != nil {
			t.Fatalf("interval %q: unexpected error: %v", tc.interval, err)
		}
		if got != tc.want {
			t.Errorf("interval %q: got %v, want %v", tc.interval, got, tc.want)
		}
	}
}

func TestDuration_CustomOverridesInterval(t *testing.T) {
	s := schedule.Schedule{
		Interval: schedule.IntervalDaily,
		Custom:   90 * time.Minute,
	}
	got, err := s.Duration()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 90*time.Minute {
		t.Errorf("got %v, want 90m", got)
	}
}

func TestDuration_UnknownInterval_ReturnsError(t *testing.T) {
	s := schedule.Schedule{Interval: "monthly"}
	_, err := s.Duration()
	if err == nil {
		t.Fatal("expected error for unknown interval, got nil")
	}
}

func TestNextRun(t *testing.T) {
	last := time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC)
	s := schedule.Schedule{Interval: schedule.IntervalHourly}
	next, err := s.NextRun(last)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := last.Add(time.Hour)
	if !next.Equal(want) {
		t.Errorf("got %v, want %v", next, want)
	}
}

func TestIsDue_WhenOverdue_ReturnsTrue(t *testing.T) {
	s := schedule.Schedule{Interval: schedule.IntervalHourly}
	last := time.Now().Add(-2 * time.Hour)
	due, err := s.IsDue(last, time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !due {
		t.Error("expected IsDue=true, got false")
	}
}

func TestIsDue_WhenNotYetDue_ReturnsFalse(t *testing.T) {
	s := schedule.Schedule{Interval: schedule.IntervalDaily}
	last := time.Now().Add(-30 * time.Minute)
	due, err := s.IsDue(last, time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if due {
		t.Error("expected IsDue=false, got true")
	}
}
