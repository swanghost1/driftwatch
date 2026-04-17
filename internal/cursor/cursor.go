// Package cursor tracks the last-processed position within a result set,
// allowing incremental runs to resume from where they left off.
package cursor

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// ErrNotFound is returned when no cursor file exists.
var ErrNotFound = errors.New("cursor: no state file found")

// State holds the cursor position for a named stream.
type State struct {
	Name      string    `json:"name"`
	Offset    int       `json:"offset"`
	LastSeen  time.Time `json:"last_seen"`
	RunCount  int       `json:"run_count"`
}

// Store persists cursor state to disk.
type Store struct {
	dir string
}

// NewStore returns a Store rooted at dir.
func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

func (s *Store) path(name string) string {
	return filepath.Join(s.dir, name+".cursor.json")
}

// Save writes state to disk, creating intermediate directories as needed.
func (s *Store) Save(st State) error {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return err
	}
	f, err := os.Create(s.path(st.Name))
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(st)
}

// Load reads the cursor state for name. Returns ErrNotFound if absent.
func (s *Store) Load(name string) (State, error) {
	f, err := os.Open(s.path(name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return State{}, ErrNotFound
		}
		return State{}, err
	}
	defer f.Close()
	var st State
	return st, json.NewDecoder(f).Decode(&st)
}

// Advance increments the offset and run count, updating LastSeen to now.
func Advance(st State, delta int) State {
	st.Offset += delta
	st.RunCount++
	st.LastSeen = time.Now().UTC()
	return st
}
