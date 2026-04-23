// Package watermark tracks high-water marks for drift counts across runs,
// recording the maximum number of drifted services ever observed.
package watermark

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Mark holds the high-water mark state.
type Mark struct {
	Peak      int       `json:"peak"`
	RecordedAt time.Time `json:"recorded_at"`
}

// Store persists and retrieves the high-water mark.
type Store struct {
	path string
}

// NewStore returns a Store backed by the given file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Load reads the current high-water mark from disk.
// Returns a zero Mark and no error if the file does not exist.
func (s *Store) Load() (Mark, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Mark{}, nil
		}
		return Mark{}, fmt.Errorf("watermark: read %s: %w", s.path, err)
	}
	var m Mark
	if err := json.Unmarshal(data, &m); err != nil {
		return Mark{}, fmt.Errorf("watermark: decode: %w", err)
	}
	return m, nil
}

// Save writes the mark to disk, creating intermediate directories as needed.
func (s *Store) Save(m Mark) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("watermark: mkdir: %w", err)
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("watermark: encode: %w", err)
	}
	return os.WriteFile(s.path, data, 0o644)
}

// Update checks the current results against the stored mark, updating it if
// the number of drifted services exceeds the previous peak. Returns the
// (possibly updated) Mark and whether a new peak was set.
func (s *Store) Update(results []drift.Result) (Mark, bool, error) {
	current, err := s.Load()
	if err != nil {
		return Mark{}, false, err
	}
	count := 0
	for _, r := range results {
		if r.Drifted {
			count++
		}
	}
	if count <= current.Peak {
		return current, false, nil
	}
	updated := Mark{Peak: count, RecordedAt: time.Now().UTC()}
	if err := s.Save(updated); err != nil {
		return Mark{}, false, err
	}
	return updated, true, nil
}

// Write prints the high-water mark as human-readable text.
func Write(w io.Writer, m Mark) {
	if m.Peak == 0 {
		fmt.Fprintln(w, "high-water mark: no drift recorded")
		return
	}
	fmt.Fprintf(w, "high-water mark: %d drifted service(s) (recorded %s)\n",
		m.Peak, m.RecordedAt.Format(time.RFC3339))
}
