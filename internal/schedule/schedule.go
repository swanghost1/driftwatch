// Package schedule provides utilities for parsing and evaluating
// cron-style run schedules for drift checks.
package schedule

import (
	"fmt"
	"time"
)

// Interval represents a named polling interval.
type Interval string

const (
	IntervalHourly  Interval = "hourly"
	IntervalDaily   Interval = "daily"
	IntervalWeekly  Interval = "weekly"
)

// Schedule holds configuration for when drift checks should run.
type Schedule struct {
	Interval Interval      `yaml:"interval"`
	Custom   time.Duration `yaml:"-"`
}

// Duration returns the time.Duration corresponding to the schedule interval.
func (s Schedule) Duration() (time.Duration, error) {
	if s.Custom > 0 {
		return s.Custom, nil
	}
	switch s.Interval {
	case IntervalHourly:
		return time.Hour, nil
	case IntervalDaily:
		return 24 * time.Hour, nil
	case IntervalWeekly:
		return 7 * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unknown interval %q", s.Interval)
	}
}

// NextRun returns the next time a check should run given the last run time.
func (s Schedule) NextRun(last time.Time) (time.Time, error) {
	d, err := s.Duration()
	if err != nil {
		return time.Time{}, err
	}
	return last.Add(d), nil
}

// IsDue reports whether a check is due relative to now, given the last run time.
func (s Schedule) IsDue(last time.Time, now time.Time) (bool, error) {
	next, err := s.NextRun(last)
	if err != nil {
		return false, err
	}
	return !now.Before(next), nil
}

// TimeUntilNext returns the duration remaining until the next scheduled run.
// If the check is already due, it returns 0.
func (s Schedule) TimeUntilNext(last time.Time, now time.Time) (time.Duration, error) {
	next, err := s.NextRun(last)
	if err != nil {
		return 0, err
	}
	if !now.Before(next) {
		return 0, nil
	}
	return next.Sub(now), nil
}
