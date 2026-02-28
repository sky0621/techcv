package sqlite

import (
	"errors"
	"os"
	"path/filepath"
)

// Prepare creates the parent directory and file for local SQLite development.
func Prepare(path string) error {
	if path == "" {
		return errors.New("sqlite path is empty")
	}

	if path == ":memory:" {
		return nil
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}

	return f.Close()
}
