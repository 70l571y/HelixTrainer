package editor

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOpenEditorReturnsMissingFileError(t *testing.T) {
	err := OpenEditor(filepath.Join(t.TempDir(), "missing.go"), "")
	if err == nil {
		t.Fatal("OpenEditor() error = nil, want missing file error")
	}
}

func TestOpenEditorReturnsHelpfulErrorWhenBinaryMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "challenge.go")
	if err := os.WriteFile(path, []byte("package main\n"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	t.Setenv("PATH", dir)

	err := OpenEditor(path, dir)
	if err == nil {
		t.Fatal("OpenEditor() error = nil, want PATH error")
	}
	if !strings.Contains(err.Error(), "hx") {
		t.Fatalf("OpenEditor() error = %v, want hx mention", err)
	}
	if !errors.Is(err, ErrHelixNotFound) {
		t.Fatalf("OpenEditor() error = %v, want ErrHelixNotFound", err)
	}
}
