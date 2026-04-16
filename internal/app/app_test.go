package app

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/70l571y/HelixTrainer/internal/challenges"
	"github.com/spf13/cobra"
)

func newTestRootCommand(t *testing.T) (*cobra.Command, *bytes.Buffer) {
	t.Helper()

	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	root := &cobra.Command{
		Use:           "hxtrainer",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	InitCommands(root, t.TempDir())

	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	return root, &out
}

func TestConfirmStatsResetPromptReadsFromIO(t *testing.T) {
	t.Run("yes", func(t *testing.T) {
		var out bytes.Buffer
		ok, err := confirmStatsReset(strings.NewReader("y\n"), &out)
		if err != nil {
			t.Fatalf("confirmStatsReset() error = %v", err)
		}
		if !ok {
			t.Fatalf("confirmStatsReset() = false, want true")
		}
		if got := strings.ToLower(out.String()); !strings.Contains(got, "будут удалены все попытки") || !strings.Contains(got, "рекорды") {
			t.Fatalf("prompt output = %q, want warning about deleting statistics", got)
		}
	})

	t.Run("no", func(t *testing.T) {
		var out bytes.Buffer
		ok, err := confirmStatsReset(strings.NewReader("n\n"), &out)
		if err != nil {
			t.Fatalf("confirmStatsReset() error = %v", err)
		}
		if ok {
			t.Fatalf("confirmStatsReset() = true, want false")
		}
	})
}

func TestStatsResetCancelsOnNo(t *testing.T) {
	oldResetAttempts := resetAttempts
	defer func() { resetAttempts = oldResetAttempts }()

	called := false
	resetAttempts = func() (int64, error) {
		called = true
		return 0, nil
	}

	root, out := newTestRootCommand(t)
	root.SetArgs([]string{"stats", "reset"})
	root.SetIn(strings.NewReader("n\n"))

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}
	if called {
		t.Fatal("resetAttempts() was called on cancel")
	}
	if got := out.String(); !strings.Contains(strings.ToLower(got), "отмен") {
		t.Fatalf("output = %q, want cancel message", got)
	}
}

func TestStatsResetReturnsErrorOnEOFWithoutConfirmation(t *testing.T) {
	oldResetAttempts := resetAttempts
	defer func() { resetAttempts = oldResetAttempts }()

	called := false
	resetAttempts = func() (int64, error) {
		called = true
		return 0, nil
	}

	root, _ := newTestRootCommand(t)
	root.SetArgs([]string{"stats", "reset"})
	root.SetIn(strings.NewReader(""))

	err := root.Execute()
	if err == nil {
		t.Fatal("root.Execute() error = nil, want error")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "подтверж") {
		t.Fatalf("root.Execute() error = %v, want confirmation error", err)
	}
	if called {
		t.Fatal("resetAttempts() was called on EOF")
	}
}

func TestStatsResetConfirmsOnYes(t *testing.T) {
	oldResetAttempts := resetAttempts
	defer func() { resetAttempts = oldResetAttempts }()

	called := 0
	resetAttempts = func() (int64, error) {
		called++
		return 3, nil
	}

	root, out := newTestRootCommand(t)
	root.SetArgs([]string{"stats", "reset"})
	root.SetIn(strings.NewReader("y\n"))

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}
	if called != 1 {
		t.Fatalf("resetAttempts() called %d times, want 1", called)
	}
	if got := out.String(); !strings.Contains(got, "3") {
		t.Fatalf("output = %q, want reset count", got)
	}
}

func TestStatsResetBypassesPromptWithYesFlag(t *testing.T) {
	oldResetAttempts := resetAttempts
	defer func() { resetAttempts = oldResetAttempts }()

	called := 0
	resetAttempts = func() (int64, error) {
		called++
		return 1, nil
	}

	root, out := newTestRootCommand(t)
	root.SetArgs([]string{"stats", "reset", "--yes"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}
	if called != 1 {
		t.Fatalf("resetAttempts() called %d times, want 1", called)
	}
	if got := out.String(); !strings.Contains(got, "1") {
		t.Fatalf("output = %q, want reset count", got)
	}
}

func TestStatsResetReturnsErrorWhenResetFails(t *testing.T) {
	oldResetAttempts := resetAttempts
	defer func() { resetAttempts = oldResetAttempts }()

	resetAttempts = func() (int64, error) {
		return 0, fmt.Errorf("reset failed")
	}

	root, _ := newTestRootCommand(t)
	root.SetArgs([]string{"stats", "reset", "--yes"})

	err := root.Execute()
	if err == nil {
		t.Fatal("root.Execute() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "reset failed") {
		t.Fatalf("root.Execute() error = %v, want reset failure", err)
	}
}

func TestStatsResetWorksTwiceOnEmptyState(t *testing.T) {
	oldResetAttempts := resetAttempts
	defer func() { resetAttempts = oldResetAttempts }()

	calls := 0
	resetAttempts = func() (int64, error) {
		calls++
		if calls == 1 {
			return 5, nil
		}
		return 0, nil
	}

	firstRoot, firstOut := newTestRootCommand(t)
	firstRoot.SetArgs([]string{"stats", "reset", "--yes"})
	if err := firstRoot.Execute(); err != nil {
		t.Fatalf("first Execute() error = %v", err)
	}

	secondRoot, secondOut := newTestRootCommand(t)
	secondRoot.SetArgs([]string{"stats", "reset", "--yes"})
	if err := secondRoot.Execute(); err != nil {
		t.Fatalf("second Execute() error = %v", err)
	}

	if calls != 2 {
		t.Fatalf("resetAttempts() called %d times, want 2", calls)
	}
	if got := firstOut.String(); !strings.Contains(got, "5") {
		t.Fatalf("first output = %q, want first reset count", got)
	}
	if got := secondOut.String(); !strings.Contains(got, "0") {
		t.Fatalf("second output = %q, want empty reset count", got)
	}
}

func TestStatsKeepsNoAttemptsBehavior(t *testing.T) {
	root, out := newTestRootCommand(t)
	root.SetArgs([]string{"stats"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}
	if got := out.String(); !strings.Contains(got, "Попыток пока нет.") {
		t.Fatalf("output = %q, want no-attempts message", got)
	}
}

func TestStatsHelpShowsResetSubcommand(t *testing.T) {
	root, out := newTestRootCommand(t)
	root.SetArgs([]string{"stats", "--help"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}
	got := strings.ToLower(out.String())
	if !strings.Contains(got, "reset") {
		t.Fatalf("help output = %q, want reset subcommand", out.String())
	}
}

func TestStatsResetHelpShowsYesFlag(t *testing.T) {
	root, out := newTestRootCommand(t)
	root.SetArgs([]string{"stats", "reset", "--help"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}
	got := out.String()
	if !strings.Contains(got, "--yes") || !strings.Contains(got, "без интерактивного подтверждения") {
		t.Fatalf("help output = %q, want reset description", got)
	}
	if !strings.Contains(got, "--yes") {
		t.Fatalf("help output = %q, want yes flag", got)
	}
}

func TestChallengeMainFileName(t *testing.T) {
	t.Run("defaults to challenge plus extension", func(t *testing.T) {
		challenge := challenges.Challenge{StartPath: "/tmp/start.go"}

		if got := challengeMainFileName(challenge); got != "challenge.go" {
			t.Fatalf("challengeMainFileName() = %q, want %q", got, "challenge.go")
		}
	})

	t.Run("uses configured main file name", func(t *testing.T) {
		challenge := challenges.Challenge{
			StartPath:    "/tmp/start.go",
			MainFileName: "matrix_processor.go",
		}

		if got := challengeMainFileName(challenge); got != "matrix_processor.go" {
			t.Fatalf("challengeMainFileName() = %q, want %q", got, "matrix_processor.go")
		}
	})
}

func TestPrepareGitWorkspaceAppliesDirtyFiles(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}

	challengeDir := t.TempDir()
	tmpDir := t.TempDir()

	challenge := challenges.Challenge{
		ID:           "git_dirty_workspace",
		Description:  "Dirty workspace",
		Language:     "go",
		DirPath:      challengeDir,
		GitDirtyFiles: map[string]string{
			"helper.go": "dirty_helper.go",
		},
	}

	targetPath := filepath.Join(tmpDir, "helper.go")
	if err := os.WriteFile(targetPath, []byte("clean\n"), 0644); err != nil {
		t.Fatalf("WriteFile(clean) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(challengeDir, "dirty_helper.go"), []byte("dirty\n"), 0644); err != nil {
		t.Fatalf("WriteFile(dirty fixture) error = %v", err)
	}

	if err := prepareGitWorkspace(challenge, tmpDir, "challenge.go"); err != nil {
		t.Fatalf("prepareGitWorkspace() error = %v", err)
	}

	content, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("ReadFile(target) error = %v", err)
	}
	if got := string(content); got != "dirty\n" {
		t.Fatalf("target content = %q, want %q", got, "dirty\n")
	}
}
