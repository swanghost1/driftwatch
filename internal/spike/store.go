package spike

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// entry is a single persisted drift count sample.
type entry struct {
	RecordedAt time.Time `json:"recorded_at"`
	DriftCount int       `json:"drift_count"`
}

// Store persists a time-ordered series of drift counts used by Detect.
type Store struct {
	path string
}

// NewStore returns a Store backed by the given file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Record appends a new drift count sample.
func (s *Store) Record(count int, at time.Time) error {
	entries, _ := s.load()
	entries = append(entries, entry{RecordedAt: at, DriftCount: count})
	return s.save(entries)
}

// Counts returns drift counts ordered oldest-first.
func (s *Store) Counts() ([]int, error) {
	entries, err := s.load()
	if err != nil {
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].RecordedAt.Before(entries[j].RecordedAt)
	})
	out := make([]int, len(entries))
	for i, e := range entries {
		out[i] = e.DriftCount
	}
	return out, nil
}

func (s *Store) load() ([]entry, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var entries []entry
	return entries, json.Unmarshal(data, &entries)
}

func (s *Store) save(entries []entry) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
