package database

import (
	"os"
	"testing"
)

func resetTestDBState(t *testing.T) {
	t.Helper()
	db = nil
}

func TestResetAttemptsClearsAttemptsAndReturnsCount(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	resetTestDBState(t)

	if err := InitDB(); err != nil {
		t.Fatalf("InitDB() error = %v", err)
	}

	if _, err := LogAttempt("alpha", true, 1.5); err != nil {
		t.Fatalf("LogAttempt(alpha) error = %v", err)
	}
	if _, err := LogAttempt("beta", false, 2.5); err != nil {
		t.Fatalf("LogAttempt(beta) error = %v", err)
	}

	removed, err := ResetAttempts()
	if err != nil {
		t.Fatalf("ResetAttempts() error = %v", err)
	}
	if removed != 2 {
		t.Fatalf("ResetAttempts() removed = %d, want 2", removed)
	}

	attempts, err := GetAllAttempts()
	if err != nil {
		t.Fatalf("GetAllAttempts() error = %v", err)
	}
	if len(attempts) != 0 {
		t.Fatalf("GetAllAttempts() len = %d, want 0", len(attempts))
	}

	removedAgain, err := ResetAttempts()
	if err != nil {
		t.Fatalf("second ResetAttempts() error = %v", err)
	}
	if removedAgain != 0 {
		t.Fatalf("second ResetAttempts() removed = %d, want 0", removedAgain)
	}
}

func TestInitDBCanRecoverAfterInitialFailure(t *testing.T) {
	resetTestDBState(t)

	blocker := t.TempDir() + "/config-file"
	if err := os.WriteFile(blocker, []byte("x"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	t.Setenv("XDG_CONFIG_HOME", blocker)

	if err := InitDB(); err == nil {
		t.Fatal("InitDB() error = nil, want initial failure")
	}

	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	if err := InitDB(); err != nil {
		t.Fatalf("second InitDB() error = %v, want recovery after failure", err)
	}

	if _, err := LogAttempt("alpha", true, 1.0); err != nil {
		t.Fatalf("LogAttempt() error after recovery = %v", err)
	}
}
