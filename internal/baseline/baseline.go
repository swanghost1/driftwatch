package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Entry represents a saved baseline of drift results for a given run.
type Entry struct {
	CreatedAt time.Time          `json:"created_at"`
	Label     string             `json:"label"`
	Results   []drift.Result     `json:"results"`
}

// Save writes a baseline entry to the given file path as JSON.
func Save(path string, label string, results []drift.Result) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("baseline: create directories: %w", err)
	}

	entry := Entry{
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Results:   results,
	}

	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("baseline: write file: %w", err)
	}

	return nil
}

// Load reads a baseline entry from the given file path.
func Load(path string) (*Entry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("baseline: file not found: %s", path)
		}
		return nil, fmt.Errorf("baseline: read file: %w", err)
	}

	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal: %w", err)
	}

	return &entry, nil
}

// Compare returns only the drift results that are new compared to the baseline.
// A result is considered new if no result with the same ServiceName and Field exists in the baseline.
func Compare(baseline *Entry, current []drift.Result) []drift.Result {
	type key struct{ service, field string }
	existing := make(map[key]struct{}, len(baseline.Results))
	for _, r := range baseline.Results {
		existing[key{r.ServiceName, r.Field}] = struct{}{}
	}

	var novel []drift.Result
	for _, r := range current {
		if _, found := existing[key{r.ServiceName, r.Field}]; !found {
			novel = append(novel, r)
		}
	}
	return novel
}
