// Package throttle guards drift check runs against being triggered too
// frequently by enforcing a configurable cooldown window between executions.
//
// Usage:
//
//	store := throttle.NewStore(".driftwatch/throttle.json")
//
//	if err := store.Check(15 * time.Minute); err != nil {
//		if errors.Is(err, throttle.ErrThrottled) {
//			fmt.Println("skipping run:", err)
//			return
//		}
//		return err
//	}
//
//	// … perform drift check …
//
//	if err := store.Record(); err != nil {
//		return err
//	}
//
// The state is stored as a small JSON file so it survives process restarts
// and can be inspected or cleared manually.
package throttle
