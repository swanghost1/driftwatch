// Package history records and retrieves past drift detection runs.
package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
)

// Entry represents a single recorded drift detection run.
type Entry struct {
	RunAt   time.Time      `json:"run_at"`
	Results []drift.Result `json:"results"`
}

// Store manages history entries persisted to disk.
type Store struct {
	Dir string
}

// NewStore returns a Store rooted at dir.
func NewStore(dir string) *Store {
	return &Store{Dir: dir}
}

// Record saves a new history entry for the current run.
func (s *Store) Record(results []drift.Result) error {
	if err := os.MkdirAll(s.Dir, 0o755); err != nil {
		return fmt.Errorf("history: create dir: %w", err)
	}
	entry := Entry{
		RunAt:   time.Now().UTC(),
		Results: results,
	}
	name := entry.RunAt.Format("20060102T150405Z") + ".json"
	path := filepath.Join(s.Dir, name)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("history: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(entry); err != nil {
		return fmt.Errorf("history: encode: %w", err)
	}
	return nil
}

// List returns all history entries sorted oldest-first.
func (s *Store) List() ([]Entry, error) {
	matches, err := filepath.Glob(filepath.Join(s.Dir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("history: glob: %w", err)
	}
	sort.Strings(matches)
	var entries []Entry
	for _, m := range matches {
		e, err := loadEntry(m)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func loadEntry(path string) (Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return Entry{}, fmt.Errorf("history: open %s: %w", path, err)
	}
	defer f.Close()
	var e Entry
	if err := json.NewDecoder(f).Decode(&e); err != nil {
		return Entry{}, fmt.Errorf("history: decode %s: %w", path, err)
	}
	return e, nil
}
