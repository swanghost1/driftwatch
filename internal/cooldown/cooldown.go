// Package cooldown tracks per-service drift events and suppresses
// repeated notifications until a configurable quiet period has elapsed.
package cooldown

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// Entry records the last time a notification was emitted for a service.
type Entry struct {
	Service  string    `json:"service"`
	LastSent time.Time `json:"last_sent"`
}

// Store holds cooldown state keyed by service name.
type Store struct {
	path    string
	entries map[string]Entry
}

// NewStore opens (or creates) a cooldown store at the given path.
func NewStore(path string) (*Store, error) {
	s := &Store{path: path, entries: make(map[string]Entry)}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	var list []Entry
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	for _, e := range list {
		s.entries[e.Service] = e
	}
	return s, nil
}

// IsCoolingDown reports whether service is still within the quiet period.
func (s *Store) IsCoolingDown(service string, period time.Duration) bool {
	e, ok := s.entries[service]
	if !ok {
		return false
	}
	return time.Since(e.LastSent) < period
}

// Record marks now as the last notification time for service and persists.
func (s *Store) Record(service string) error {
	s.entries[service] = Entry{Service: service, LastSent: time.Now()}
	return s.save()
}

func (s *Store) save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	list := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		list = append(list, e)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
