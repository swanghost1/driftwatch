package variance

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Entry is a persisted snapshot of one run's results, used to reconstruct
// historical drift counts for variance analysis.
type Entry struct {
	RecordedAt time.Time     `json:"recorded_at"`
	Results    []drift.Result `json:"results"`
}

// Store persists and retrieves variance history entries on disk.
type Store struct {
	dir string
}

// NewStore returns a Store rooted at dir.
func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

// Record appends a new entry for the current run.
func (s *Store) Record(results []drift.Result) error {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return fmt.Errorf("variance store: mkdir: %w", err)
	}
	entry := Entry{
		RecordedAt: time.Now().UTC(),
		Results:    results,
	}
	name := entry.RecordedAt.Format("20060102T150405.000000000Z") + ".json"
	path := filepath.Join(s.dir, name)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("variance store: create: %w", err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(entry); err != nil {
		return fmt.Errorf("variance store: encode: %w", err)
	}
	return nil
}

// Load reads up to limit entries from the store, oldest first.
// Pass limit <= 0 to load all entries.
func (s *Store) Load(limit int) ([]Entry, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("variance store: readdir: %w", err)
	}
	var out []Entry
	for _, de := range entries {
		if de.IsDir() {
			continue
		}
		path := filepath.Join(s.dir, de.Name())
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("variance store: open %s: %w", de.Name(), err)
		}
		var e Entry
		if err := json.NewDecoder(f).Decode(&e); err != nil {
			f.Close()
			return nil, fmt.Errorf("variance store: decode %s: %w", de.Name(), err)
		}
		f.Close()
		out = append(out, e)
	}
	if limit > 0 && len(out) > limit {
		out = out[len(out)-limit:]
	}
	return out, nil
}
