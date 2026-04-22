package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootHelpShowsStatsResetExamples(t *testing.T) {
	rootCmd := newRootCommand(t.TempDir())

	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"--help"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("rootCmd.Execute() error = %v", err)
	}

	help := out.String()
	if !strings.Contains(help, "hxtrainer stats reset") {
		t.Fatalf("root help = %q, want stats reset example", help)
	}
	if !strings.Contains(help, "hxtrainer stats reset --yes") {
		t.Fatalf("root help = %q, want stats reset --yes example", help)
	}
	if !strings.Contains(help, "hxtrainer doctor") {
		t.Fatalf("root help = %q, want doctor example", help)
	}
	if !strings.Contains(help, "hxtrainer completion bash") {
		t.Fatalf("root help = %q, want completion example", help)
	}
}
