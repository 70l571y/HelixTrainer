package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/70l571y/HelixTrainer/internal/challenges"
	"github.com/70l571y/HelixTrainer/internal/database"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func newTestRootCommand(t *testing.T) (*cobra.Command, *bytes.Buffer) {
	t.Helper()

	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	return newTestRootCommandWithChallengesDir(t, t.TempDir())
}

func newTestRootCommandWithChallengesDir(t *testing.T, challengesDir string) (*cobra.Command, *bytes.Buffer) {
	t.Helper()

	root := &cobra.Command{
		Use:           "hxtrainer",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	InitCommands(root, challengesDir)
	resetCommandState(root)

	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	return root, &out
}

func resetCommandState(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		defValue := flag.DefValue
		if flag.Value.Type() == "stringSlice" && defValue == "[]" {
			defValue = ""
		}
		_ = cmd.Flags().Set(flag.Name, defValue)
		flag.Changed = false
	})
	cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		defValue := flag.DefValue
		if flag.Value.Type() == "stringSlice" && defValue == "[]" {
			defValue = ""
		}
		_ = cmd.PersistentFlags().Set(flag.Name, defValue)
		flag.Changed = false
	})
	for _, child := range cmd.Commands() {
		resetCommandState(child)
	}
}

func writeTestChallenge(t *testing.T, root string, id string) {
	t.Helper()

	dir := filepath.Join(root, id)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	config := `{
  "id": "` + id + `",
  "title": "Test Challenge",
  "description": "Сделайте правку",
  "difficulty": "Easy",
  "language": "go",
  "judge_mode": "exact",
  "start_file": "start.go",
  "goal_file": "goal.go",
  "tags": ["edit_insert"],
  "author_time": 5
}`
	if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte(config), 0644); err != nil {
		t.Fatalf("WriteFile(config) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "start.go"), []byte("//go:build ignore\n\npackage main\n"), 0644); err != nil {
		t.Fatalf("WriteFile(start) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "goal.go"), []byte("//go:build ignore\n\npackage main\n"), 0644); err != nil {
		t.Fatalf("WriteFile(goal) error = %v", err)
	}
}

func updateChallengeTags(t *testing.T, path string, tags []string) {
	t.Helper()
	updateChallengeConfigField(t, path, "tags", tags)
}

func updateChallengeDifficulty(t *testing.T, path string, difficulty string) {
	t.Helper()
	updateChallengeConfigField(t, path, "difficulty", difficulty)
}

func updateChallengeConfigField(t *testing.T, path string, field string, value any) {
	t.Helper()

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(config) error = %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("json.Unmarshal(config) error = %v", err)
	}
	payload[field] = value

	updated, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		t.Fatalf("json.MarshalIndent(config) error = %v", err)
	}
	if err := os.WriteFile(path, updated, 0644); err != nil {
		t.Fatalf("WriteFile(config) error = %v", err)
	}
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
		ID:          "git_dirty_workspace",
		Description: "Dirty workspace",
		Language:    "go",
		DirPath:     challengeDir,
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

func TestParsePostChallengeAction(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    postChallengeAction
		wantErr bool
	}{
		{name: "next", input: "j\n", want: actionNext},
		{name: "retry", input: "k\n", want: actionRetry},
		{name: "quit", input: "q\n", want: actionQuit},
		{name: "upper", input: "J\n", want: actionNext},
		{name: "invalid", input: "??\n", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePostChallengeAction(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parsePostChallengeAction() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Fatalf("parsePostChallengeAction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadPostChallengeActionRepromptsUntilValid(t *testing.T) {
	var out bytes.Buffer
	got := readPostChallengeAction(strings.NewReader("???\nk\n"), &out)
	if got != actionRetry {
		t.Fatalf("readPostChallengeAction() = %v, want %v", got, actionRetry)
	}

	output := out.String()
	if !strings.Contains(output, "Неизвестная команда") {
		t.Fatalf("output = %q, want invalid input warning", output)
	}
	if strings.Count(output, "Введите действие") != 2 {
		t.Fatalf("output = %q, want prompt twice", output)
	}
}

func TestCompareReleaseVersions(t *testing.T) {
	if got := compareReleaseVersions("1.2.0", "1.1.9"); got <= 0 {
		t.Fatalf("compareReleaseVersions() = %d, want > 0", got)
	}
	if got := compareReleaseVersions("1.2.0", "1.2.0"); got != 0 {
		t.Fatalf("compareReleaseVersions() = %d, want 0", got)
	}
	if got := compareReleaseVersions("1.1.9", "1.2.0"); got >= 0 {
		t.Fatalf("compareReleaseVersions() = %d, want < 0", got)
	}
}

func TestLatestVersionFromResponse(t *testing.T) {
	body := `{"tag_name":"v1.4.2"}`
	got, err := latestVersionFromResponse(strings.NewReader(body))
	if err != nil {
		t.Fatalf("latestVersionFromResponse() error = %v", err)
	}
	if got != "1.4.2" {
		t.Fatalf("latestVersionFromResponse() = %q, want %q", got, "1.4.2")
	}
}

func TestLatestVersionFromResponseRejectsEmptyTag(t *testing.T) {
	_, err := latestVersionFromResponse(strings.NewReader(`{"tag_name":" "}`))
	if err == nil {
		t.Fatal("latestVersionFromResponse() error = nil, want error")
	}
}

func TestGetLatestVersionHandlesResponseScenarios(t *testing.T) {
	origClient := releaseHTTPClient
	defer func() { releaseHTTPClient = origClient }()

	t.Run("ok", func(t *testing.T) {
		releaseHTTPClient = testHTTPClient(func(req *http.Request) (*http.Response, error) {
			if req.URL.String() != githubLatestReleaseURL {
				t.Fatalf("request URL = %q, want %q", req.URL.String(), githubLatestReleaseURL)
			}
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(`{"tag_name":"v2.3.4"}`)),
				Header:     make(http.Header),
			}, nil
		})

		if got := getLatestVersion(); got != "2.3.4" {
			t.Fatalf("getLatestVersion() = %q, want %q", got, "2.3.4")
		}
	})

	t.Run("request failure", func(t *testing.T) {
		releaseHTTPClient = testHTTPClient(func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("boom")
		})

		if got := getLatestVersion(); got != "" {
			t.Fatalf("getLatestVersion() = %q, want empty string", got)
		}
	})

	t.Run("bad status", func(t *testing.T) {
		releaseHTTPClient = testHTTPClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 503,
				Body:       io.NopCloser(strings.NewReader(`{"tag_name":"v2.3.4"}`)),
				Header:     make(http.Header),
			}, nil
		})

		if got := getLatestVersion(); got != "" {
			t.Fatalf("getLatestVersion() = %q, want empty string", got)
		}
	})
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func testHTTPClient(fn roundTripFunc) *http.Client {
	return &http.Client{Transport: fn}
}

func TestListJSONOutput(t *testing.T) {
	challengesDir := t.TempDir()
	writeTestChallenge(t, challengesDir, "json_list_case")

	root, out := newTestRootCommandWithChallengesDir(t, challengesDir)
	root.SetArgs([]string{"list", "--json"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	var payload []map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal() error = %v; output=%q", err, out.String())
	}
	if len(payload) != 1 {
		t.Fatalf("payload len = %d, want 1", len(payload))
	}
	if payload[0]["id"] != "json_list_case" {
		t.Fatalf("payload[0][id] = %v, want %q", payload[0]["id"], "json_list_case")
	}
}

func TestListJSONOutputWithTrackFilter(t *testing.T) {
	challengesDir := t.TempDir()
	writeTestChallenge(t, challengesDir, "core_case")
	writeTestChallenge(t, challengesDir, "optional_case")

	updateChallengeTags(t, filepath.Join(challengesDir, "core_case", "config.json"), []string{"movement_basic", "track_core_hotkey"})
	updateChallengeTags(t, filepath.Join(challengesDir, "optional_case", "config.json"), []string{"command_mode", "track_optional_command_line"})

	root, out := newTestRootCommandWithChallengesDir(t, challengesDir)
	root.SetArgs([]string{"list", "--json", "--track", "core"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	var payload []map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal() error = %v; output=%q", err, out.String())
	}
	if len(payload) != 1 || payload[0]["id"] != "core_case" {
		t.Fatalf("payload = %#v, want only core_case", payload)
	}
}

func TestStatsJSONOutputWithoutAttempts(t *testing.T) {
	challengesDir := t.TempDir()
	writeTestChallenge(t, challengesDir, "json_stats_case")

	root, out := newTestRootCommandWithChallengesDir(t, challengesDir)
	root.SetArgs([]string{"stats", "--json"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	var payload struct {
		RecentAttempts []map[string]any `json:"recent_attempts"`
		Challenges     []map[string]any `json:"challenges"`
	}
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal() error = %v; output=%q", err, out.String())
	}
	if len(payload.RecentAttempts) != 0 {
		t.Fatalf("recent_attempts len = %d, want 0", len(payload.RecentAttempts))
	}
	if len(payload.Challenges) != 1 {
		t.Fatalf("challenges len = %d, want 1", len(payload.Challenges))
	}
	if payload.Challenges[0]["id"] != "json_stats_case" {
		t.Fatalf("challenges[0][id] = %v, want %q", payload.Challenges[0]["id"], "json_stats_case")
	}
}

func TestStatsJSONOutputWithDifficultyFilter(t *testing.T) {
	challengesDir := t.TempDir()
	writeTestChallenge(t, challengesDir, "easy_case")
	writeTestChallenge(t, challengesDir, "medium_case")

	updateChallengeDifficulty(t, filepath.Join(challengesDir, "medium_case", "config.json"), "Medium")

	root, out := newTestRootCommandWithChallengesDir(t, challengesDir)
	root.SetArgs([]string{"stats", "--json", "--difficulty", "medium"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	var payload struct {
		Challenges []map[string]any `json:"challenges"`
	}
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal() error = %v; output=%q", err, out.String())
	}
	if len(payload.Challenges) != 1 || payload.Challenges[0]["id"] != "medium_case" {
		t.Fatalf("payload = %#v, want only medium_case", payload.Challenges)
	}
}

func TestCompletionCommandHelpAndBashOutput(t *testing.T) {
	root, out := newTestRootCommand(t)
	root.SetArgs([]string{"completion", "bash"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}
	if got := out.String(); !strings.Contains(got, "hxtrainer") {
		t.Fatalf("completion output = %q, want shell completion script", got)
	}
}

func TestDoctorReport(t *testing.T) {
	origLookPath := binaryLookPath
	defer func() { binaryLookPath = origLookPath }()

	binaryLookPath = func(name string) (string, error) {
		if name == "hx" {
			return "", fmt.Errorf("missing")
		}
		return "/usr/bin/" + name, nil
	}

	checks := collectDoctorChecks("/tmp/config", "/tmp/db.sqlite", "/tmp/challenges")
	if len(checks) < 4 {
		t.Fatalf("collectDoctorChecks() len = %d, want at least 4", len(checks))
	}

	var foundHelix bool
	for _, check := range checks {
		if check.Name == "helix" {
			foundHelix = true
			if check.OK {
				t.Fatalf("helix check = %+v, want failed status", check)
			}
		}
	}
	if !foundHelix {
		t.Fatal("collectDoctorChecks() missing helix check")
	}
}

func TestDoctorJSONOutput(t *testing.T) {
	origLookPath := binaryLookPath
	defer func() { binaryLookPath = origLookPath }()

	binaryLookPath = func(name string) (string, error) {
		return "/usr/bin/" + name, nil
	}

	root, out := newTestRootCommand(t)
	root.SetArgs([]string{"doctor", "--json"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	var payload []map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal() error = %v; output=%q", err, out.String())
	}
	if len(payload) == 0 {
		t.Fatal("doctor json payload is empty")
	}
}

func TestSelectChallengeForStrategy(t *testing.T) {
	challengeList := []challenges.Challenge{
		{ID: "easy_done", Difficulty: "Easy", Tags: []string{"movement_basic", "track_core_hotkey"}},
		{ID: "easy_next", Difficulty: "Easy", Tags: []string{"movement_basic", "track_core_hotkey"}},
		{ID: "lsp_new", Difficulty: "Easy", Tags: []string{"lsp_reference", "track_core_hotkey"}},
	}
	attempts := []database.Attempt{
		{ChallengeID: "easy_done", IsCorrect: true},
		{ChallengeID: "easy_next", IsCorrect: false},
	}

	got, err := selectChallengeForStrategy(challengeList, "", "progression", attempts)
	if err != nil {
		t.Fatalf("selectChallengeForStrategy() error = %v", err)
	}
	if got.ID != "easy_next" {
		t.Fatalf("selectChallengeForStrategy(progression) = %q, want %q", got.ID, "easy_next")
	}

	got, err = selectChallengeForStrategy(challengeList, "", "weak-skills", attempts)
	if err != nil {
		t.Fatalf("selectChallengeForStrategy() error = %v", err)
	}
	if got.ID != "lsp_new" {
		t.Fatalf("selectChallengeForStrategy(weak-skills) = %q, want %q", got.ID, "lsp_new")
	}
}

func TestQueueJSONOutput(t *testing.T) {
	challengesDir := t.TempDir()
	writeTestChallenge(t, challengesDir, "queue_one")
	writeTestChallenge(t, challengesDir, "queue_two")
	updateChallengeDifficulty(t, filepath.Join(challengesDir, "queue_two", "config.json"), "Medium")

	root, out := newTestRootCommandWithChallengesDir(t, challengesDir)
	root.SetArgs([]string{"queue", "--json", "--limit", "1", "--strategy", "progression"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	var payload []map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal() error = %v; output=%q", err, out.String())
	}
	if len(payload) != 1 {
		t.Fatalf("payload len = %d, want 1", len(payload))
	}
	if payload[0]["id"] != "queue_one" {
		t.Fatalf("payload[0][id] = %v, want %q", payload[0]["id"], "queue_one")
	}
}

func TestHistoryJSONOutput(t *testing.T) {
	challengesDir := t.TempDir()
	writeTestChallenge(t, challengesDir, "history_case")

	root, _ := newTestRootCommandWithChallengesDir(t, challengesDir)
	if _, err := database.LogAttempt("history_case", true, 4.2); err != nil {
		t.Fatalf("LogAttempt() error = %v", err)
	}

	root, out := newTestRootCommandWithChallengesDir(t, challengesDir)
	root.SetArgs([]string{"history", "history_case", "--json"})
	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	var payload []map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal() error = %v; output=%q", err, out.String())
	}
	if len(payload) == 0 {
		t.Fatal("history payload is empty")
	}
	if payload[0]["challenge_id"] != "history_case" {
		t.Fatalf("payload[0][challenge_id] = %v, want %q", payload[0]["challenge_id"], "history_case")
	}
}

func TestStatsExportImportRoundTrip(t *testing.T) {
	challengesDir := t.TempDir()
	writeTestChallenge(t, challengesDir, "export_case")

	exportPath := filepath.Join(t.TempDir(), "attempts.json")

	root, _ := newTestRootCommandWithChallengesDir(t, challengesDir)
	if _, err := database.LogAttempt("export_case", true, 2.5); err != nil {
		t.Fatalf("LogAttempt() error = %v", err)
	}

	root, _ = newTestRootCommandWithChallengesDir(t, challengesDir)
	root.SetArgs([]string{"stats", "export", exportPath})
	if err := root.Execute(); err != nil {
		t.Fatalf("stats export error = %v", err)
	}

	resetTestDB := func() {
		if _, err := database.ResetAttempts(); err != nil {
			t.Fatalf("ResetAttempts() error = %v", err)
		}
	}
	resetTestDB()

	root, _ = newTestRootCommandWithChallengesDir(t, challengesDir)
	root.SetArgs([]string{"stats", "import", exportPath})
	if err := root.Execute(); err != nil {
		t.Fatalf("stats import error = %v", err)
	}

	attempts, err := database.GetAttemptsByChallenge("export_case")
	if err != nil {
		t.Fatalf("GetAttemptsByChallenge() error = %v", err)
	}
	if len(attempts) != 1 {
		t.Fatalf("attempts len = %d, want 1", len(attempts))
	}
}
