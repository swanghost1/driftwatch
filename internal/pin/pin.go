// Package pin allows specific drift results to be "pinned" — marked as
// acknowledged and intentionally allowed to differ from declared state.
// Pinned results are excluded from drift counts and reports.
package pin

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Rule describes a single pin entry.
type Rule struct {
	Service string    `json:"service"`
	Field   string    `json:"field"`
	Reason  string    `json:"reason"`
	PinnedAt time.Time `json:"pinned_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// Store holds a collection of pin rules persisted to disk.
type Store struct {
	path  string
	Rules []Rule `json:"rules"`
}

// LoadStore reads a pin store from path. If the file does not exist an empty
// store is returned without error.
func LoadStore(path string) (*Store, error) {
	s := &Store{path: path}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, s); err != nil {
		return nil, err
	}
	return s, nil
}

// SaveStore writes the store back to its original path.
func SaveStore(s *Store) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

// Apply returns a copy of results with pinned drift entries marked as not
// drifted and annotated with the pin reason.
func Apply(results []drift.Result, s *Store) []drift.Result {
	out := make([]drift.Result, 0, len(results))
	for _, r := range results {
		if isPinned(r, s, time.Now()) {
			r.Drifted = false
			r.Live = "[pinned] " + r.Live
		}
		out = append(out, r)
	}
	return out
}

func isPinned(r drift.Result, s *Store, now time.Time) bool {
	if !r.Drifted {
		return false
	}
	for _, rule := range s.Rules {
		if rule.ExpiresAt != nil && now.After(*rule.ExpiresAt) {
			continue
		}
		svcMatch := strings.EqualFold(rule.Service, r.Service) || rule.Service == "*"
		fieldMatch := rule.Field == "" || strings.EqualFold(rule.Field, r.Field)
		if svcMatch && fieldMatch {
			return true
		}
	}
	return false
}
