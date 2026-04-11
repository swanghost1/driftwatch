// Package output handles writing drift results to various destinations.
package output

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Destination represents where output should be written.
type Destination struct {
	// Path is the file path to write to. If empty, stdout is used.
	Path string
	// Append controls whether to append to an existing file.
	Append bool
}

// Open returns a WriteCloser for the given destination.
// If Path is empty, os.Stdout is returned with a no-op closer.
func Open(d Destination) (io.WriteCloser, error) {
	if d.Path == "" {
		return nopCloser{os.Stdout}, nil
	}

	dir := filepath.Dir(d.Path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("output: create directories for %q: %w", d.Path, err)
	}

	flags := os.O_CREATE | os.O_WRONLY
	if d.Append {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}

	f, err := os.OpenFile(d.Path, flags, 0o644)
	if err != nil {
		return nil, fmt.Errorf("output: open %q: %w", d.Path, err)
	}
	return f, nil
}

// nopCloser wraps a Writer with a no-op Close method.
type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }
