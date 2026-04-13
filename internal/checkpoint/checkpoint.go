// Package checkpoint tracks the last successful run time and result summary
// for each named scan, allowing driftwatch to skip redundant work and surface
// how recently each service was evaluated.
package checkpoint

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// Entry records the outcome of a single completed scan.
type Entry struct {
	Name      string    `json:"name"`
	RunAt     time.Time `json:"run_at"`
	Total     int       `json:"total"`
	Drifted   int       `json:"drifted"`
	Clean     int       `json:"clean"`
}

// Store persists checkpoint entries on disk, one file per named scan.
type Store struct {
	dir string
}

// NewStore returns a Store rooted at dir.
func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

// Save writes e to disk, overwriting any previous checkpoint for e.Name.
func (s *Store) Save(e Entry) error {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath(e.Name), data, 0o644)
}

// Load retrieves the checkpoint for name. Returns ErrNotFound when no
// checkpoint has been saved yet.
func (s *Store) Load(name string) (Entry, error) {
	data, err := os.ReadFile(s.filePath(name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Entry{}, ErrNotFound
		}
		return Entry{}, err
	}
	var e Entry
	if err := json.Unmarshal(data, &e); err != nil {
		return Entry{}, err
	}
	return e, nil
}

// ErrNotFound is returned by Load when no checkpoint exists for the given name.
var ErrNotFound = errors.New("checkpoint: no entry found")

func (s *Store) filePath(name string) string {
	return filepath.Join(s.dir, name+".json")
}
