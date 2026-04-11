// Package schedule provides types and helpers for determining when drift
// checks should be executed based on a configured polling interval.
//
// A Schedule can be constructed with one of the predefined named intervals
// (hourly, daily, weekly) or with a custom time.Duration for finer control.
//
// Example usage:
//
//	s := schedule.Schedule{Interval: schedule.IntervalHourly}
//	due, err := s.IsDue(lastRun, time.Now())
//	if err != nil {
//		log.Fatal(err)
//	}
//	if due {
//		// run drift check
//	}
package schedule
