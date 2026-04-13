// Package mute provides time-bounded silencing of drift results for
// specific services and fields. A muted result is excluded from reports
// and notifications until the mute window expires.
package mute

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/example/driftwatch/internal/drift"
)

// Rule silences drift for a specific service and optional field.
type Rule struct {
	Service   string    `json:"service"`
	Field     string    `json:"field,omitempty"` // empty means all fields
	ExpiresAt time.Time `json:"expires_at"`
	Reason    string    `json:"reason,omitempty"`
}

// Store holds the active mute rules.
type Store struct {
	Rules []Rule `json:"rules"`
}

// LoadStore reads a mute store from path. If the file does not exist an
// empty store is returned without error.
func LoadStore(path string) (*Store, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Store{}, nil
	}
	if err != nil {
		return nil, err
	}
	var s Store
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// SaveStore persists s to path, creating intermediate directories as needed.
func SaveStore(path string, s *Store) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Apply filters results, removing any drift entry that is currently muted.
// Results with no active mute rule are returned unchanged.
func Apply(results []drift.Result, s *Store, now time.Time) []drift.Result {
	if s == nil || len(s.Rules) == 0 {
		return results
	}
	out := make([]drift.Result, 0, len(results))
	for _, r := range results {
		if !isMuted(r, s.Rules, now) {
			out = append(out, r)
		}
	}
	return out
}

func isMuted(r drift.Result, rules []Rule, now time.Time) bool {
	for _, rule := range rules {
		if now.After(rule.ExpiresAt) {
			continue
		}
		if !strings.EqualFold(rule.Service, r.Service) {
			continue
		}
		if rule.Field == "" || strings.EqualFold(rule.Field, r.Field) {
			return true
		}
	}
	return false
}
