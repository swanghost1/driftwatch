package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ServiceSnapshot represents the captured state of a service at a point in time.
type ServiceSnapshot struct {
	Name      string            `json:"name"`
	Image     string            `json:"image"`
	Replicas  int               `json:"replicas"`
	Env       map[string]string `json:"env,omitempty"`
	CapturedAt time.Time        `json:"captured_at"`
}

// Store holds a collection of service snapshots.
type Store struct {
	Version   string            `json:"version"`
	Snapshots []ServiceSnapshot `json:"snapshots"`
}

// Save writes the snapshot store to a JSON file at the given path.
func Save(path string, store *Store) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("snapshot: create directory: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(store); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// Load reads a snapshot store from a JSON file at the given path.
func Load(path string) (*Store, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("snapshot: file not found: %s", path)
		}
		return nil, fmt.Errorf("snapshot: open file: %w", err)
	}
	defer f.Close()

	var store Store
	if err := json.NewDecoder(f).Decode(&store); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	return &store, nil
}

// FindByName returns the snapshot for the named service, or nil if not found.
func (s *Store) FindByName(name string) *ServiceSnapshot {
	for i := range s.Snapshots {
		if s.Snapshots[i].Name == name {
			return &s.Snapshots[i]
		}
	}
	return nil
}
