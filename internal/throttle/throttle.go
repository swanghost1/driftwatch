// Package throttle provides rate-limiting for drift check runs,
// preventing excessive checks within a configured cooldown window.
package throttle

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ErrThrottled is returned when a run is attempted before the cooldown has elapsed.
var ErrThrottled = errors.New("throttle: cooldown period has not elapsed")

// State records the last time a drift check was successfully run.
type State struct {
	LastRun time.Time `json:"last_run"`
}

// Store persists throttle state to disk.
type Store struct {
	path string
}

// NewStore returns a Store backed by the given file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Check returns ErrThrottled if the cooldown window has not elapsed since the
// last recorded run. If no state file exists the check passes.
func (s *Store) Check(cooldown time.Duration) error {
	st, err := s.load()
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("throttle: load state: %w", err)
	}
	if elapsed := time.Since(st.LastRun); elapsed < cooldown {
		return fmt.Errorf("%w: next run allowed in %s", ErrThrottled,
			(cooldown - elapsed).Round(time.Second))
	}
	return nil
}

// Record persists the current time as the last run timestamp.
func (s *Store) Record() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("throttle: mkdir: %w", err)
	}
	data, err := json.Marshal(State{LastRun: time.Now().UTC()})
	if err != nil {
		return fmt.Errorf("throttle: marshal: %w", err)
	}
	if err := os.WriteFile(s.path, data, 0o644); err != nil {
		return fmt.Errorf("throttle: write: %w", err)
	}
	return nil
}

func (s *Store) load() (State, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return State{}, err
	}
	var st State
	if err := json.Unmarshal(data, &st); err != nil {
		return State{}, fmt.Errorf("unmarshal: %w", err)
	}
	return st, nil
}
