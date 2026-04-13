package stale

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// Store persists the firstSeen map to disk so staleness state survives
// between driftwatch invocations.
type Store struct {
	path string
}

// NewStore returns a Store that reads and writes from path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Load reads the firstSeen map from disk. If the file does not exist an
// empty (non-nil) map is returned without error.
func (s *Store) Load() (map[string]time.Time, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return make(map[string]time.Time), nil
		}
		return nil, err
	}
	var m map[string]time.Time
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// Save writes the firstSeen map to disk, creating intermediate directories
// as needed.
func (s *Store) Save(m map[string]time.Time) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
