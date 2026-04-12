// Package audit records drift check events for compliance and traceability.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/example/driftwatch/internal/drift"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp  time.Time          `json:"timestamp"`
	ConfigFile string             `json:"config_file"`
	TriggeredBy string            `json:"triggered_by"`
	TotalChecked int              `json:"total_checked"`
	DriftedCount int              `json:"drifted_count"`
	Results    []drift.Result     `json:"results"`
}

// Logger writes audit entries to a directory as newline-delimited JSON files.
type Logger struct {
	dir string
}

// NewLogger creates a Logger that stores entries under dir.
func NewLogger(dir string) *Logger {
	return &Logger{dir: dir}
}

// Record writes an audit entry for the given drift results.
func (l *Logger) Record(configFile, triggeredBy string, results []drift.Result) error {
	if err := os.MkdirAll(l.dir, 0o755); err != nil {
		return fmt.Errorf("audit: create dir: %w", err)
	}

	drifted := 0
	for _, r := range results {
		if r.Drifted {
			drifted++
		}
	}

	entry := Entry{
		Timestamp:    time.Now().UTC(),
		ConfigFile:   configFile,
		TriggeredBy:  triggeredBy,
		TotalChecked: len(results),
		DriftedCount: drifted,
		Results:      results,
	}

	filename := fmt.Sprintf("%s.json", entry.Timestamp.Format("20060102T150405Z"))
	path := filepath.Join(l.dir, filename)

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("audit: create file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(entry); err != nil {
		return fmt.Errorf("audit: encode entry: %w", err)
	}
	return nil
}

// List returns all audit entries from the directory, oldest first.
func (l *Logger) List() ([]Entry, error) {
	glob := filepath.Join(l.dir, "*.json")
	matches, err := filepath.Glob(glob)
	if err != nil {
		return nil, fmt.Errorf("audit: glob: %w", err)
	}

	var entries []Entry
	for _, m := range matches {
		data, err := os.ReadFile(m)
		if err != nil {
			return nil, fmt.Errorf("audit: read %s: %w", m, err)
		}
		var e Entry
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("audit: decode %s: %w", m, err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
