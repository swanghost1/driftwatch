package retention

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Policy defines how long entries are retained.
type Policy struct {
	// MaxAge is the maximum age of an entry before it is pruned.
	MaxAge time.Duration
	// MaxEntries is the maximum number of entries to keep (0 = unlimited).
	MaxEntries int
}

// DefaultPolicy returns a sensible default retention policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxAge:     30 * 24 * time.Hour, // 30 days
		MaxEntries: 100,
	}
}

// Entry represents a file-based record with a timestamp embedded in its name.
type Entry struct {
	Path    string
	Created time.Time
}

// Prune removes entries from dir that violate the given policy.
// Files must follow the naming pattern produced by history/audit (RFC3339 prefix).
// Returns the paths of deleted files and any error encountered.
func Prune(dir string, p Policy) ([]string, error) {
	entries, err := collect(dir)
	if err != nil {
		return nil, fmt.Errorf("retention: collect entries: %w", err)
	}

	// Sort oldest-first so we can trim by count from the front.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Created.Before(entries[j].Created)
	})

	now := time.Now()
	var pruned []string

	for _, e := range entries {
		if p.MaxAge > 0 && now.Sub(e.Created) > p.MaxAge {
			if err := os.Remove(e.Path); err != nil && !os.IsNotExist(err) {
				return pruned, fmt.Errorf("retention: remove %s: %w", e.Path, err)
			}
			pruned = append(pruned, e.Path)
		}
	}

	// Re-collect surviving entries for count-based pruning.
	survivors, err := collect(dir)
	if err != nil {
		return pruned, fmt.Errorf("retention: re-collect after age prune: %w", err)
	}
	sort.Slice(survivors, func(i, j int) bool {
		return survivors[i].Created.Before(survivors[j].Created)
	})

	if p.MaxEntries > 0 && len(survivors) > p.MaxEntries {
		excess := survivors[:len(survivors)-p.MaxEntries]
		for _, e := range excess {
			if err := os.Remove(e.Path); err != nil && !os.IsNotExist(err) {
				return pruned, fmt.Errorf("retention: remove %s: %w", e.Path, err)
			}
			pruned = append(pruned, e.Path)
		}
	}

	return pruned, nil
}

func collect(dir string) ([]Entry, error) {
	glob := filepath.Join(dir, "*.json")
	matches, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}
	var entries []Entry
	for _, m := range matches {
		base := filepath.Base(m)
		// Expect filenames like "2006-01-02T15:04:05Z07:00_*.json"
		ts, err := time.Parse(time.RFC3339, base[:len(time.RFC3339)])
		if err != nil {
			// Fall back to file mtime.
			info, serr := os.Stat(m)
			if serr != nil {
				continue
			}
			ts = info.ModTime()
		}
		entries = append(entries, Entry{Path: m, Created: ts})
	}
	return entries, nil
}
