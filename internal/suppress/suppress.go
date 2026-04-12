// Package suppress provides functionality for suppressing known or accepted
// drift results so they are excluded from reports and policy evaluation.
package suppress

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/example/driftwatch/internal/drift"
)

// Rule defines a suppression rule for a specific service and field.
type Rule struct {
	Service string    `json:"service"`
	Field   string    `json:"field"`
	Reason  string    `json:"reason"`
	Expires time.Time `json:"expires,omitempty"`
}

// IsExpired reports whether the rule has passed its expiry time.
func (r Rule) IsExpired() bool {
	return !r.Expires.IsZero() && time.Now().After(r.Expires)
}

// Store holds a collection of suppression rules.
type Store struct {
	Rules []Rule `json:"rules"`
}

// LoadStore reads suppression rules from a JSON file at the given path.
// If the file does not exist, an empty Store is returned.
func LoadStore(path string) (*Store, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Store{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("suppress: read file: %w", err)
	}
	var s Store
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("suppress: parse file: %w", err)
	}
	return &s, nil
}

// SaveStore writes the Store to a JSON file at the given path,
// creating intermediate directories as needed.
func SaveStore(path string, s *Store) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("suppress: mkdir: %w", err)
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("suppress: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("suppress: write file: %w", err)
	}
	return nil
}

// Apply filters out drift results that are matched by an active (non-expired)
// suppression rule. A rule matches when both Service and Field are non-empty
// substrings (case-insensitive) of the result's ServiceName and Field.
func Apply(results []drift.Result, s *Store) []drift.Result {
	if s == nil || len(s.Rules) == 0 {
		return results
	}
	filtered := make([]drift.Result, 0, len(results))
	for _, r := range results {
		if !isSuppressed(r, s) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

func isSuppressed(r drift.Result, s *Store) bool {
	for _, rule := range s.Rules {
		if rule.IsExpired() {
			continue
		}
		svcMatch := strings.Contains(
			strings.ToLower(r.ServiceName),
			strings.ToLower(rule.Service),
		)
		fieldMatch := strings.Contains(
			strings.ToLower(r.Field),
			strings.ToLower(rule.Field),
		)
		if svcMatch && fieldMatch {
			return true
		}
	}
	return false
}
