package challengesdata

import (
	"bytes"
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

// Embedded содержит встроенные данные челленджей.
//
//go:embed go/**
var Embedded embed.FS

// SyncToDir раскладывает встроенные данные челленджей в целевую директорию.
func SyncToDir(dstRoot string) error {
	return fs.WalkDir(Embedded, "go", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dstRoot, filepath.FromSlash(path))
		if d.IsDir() {
			return os.MkdirAll(dstPath, 0755)
		}

		data, err := Embedded.ReadFile(path)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return err
		}

		if current, err := os.ReadFile(dstPath); err == nil && bytes.Equal(current, data) {
			return nil
		}

		return os.WriteFile(dstPath, data, 0644)
	})
}
