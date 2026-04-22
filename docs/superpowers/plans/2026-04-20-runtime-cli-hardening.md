# Runtime CLI Hardening Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make HelixTrainer's CLI runtime reliable by fixing config directory handling, editor error reporting, play-loop interaction, and upgrade checking without adding new product features.

**Architecture:** Keep the current CLI shape, but extract small helper functions around config path resolution, post-challenge action parsing, and version checks so behavior is testable. Improve orchestration in `internal/app/app.go` while keeping the user-visible command set unchanged.

**Tech Stack:** Go 1.26, Cobra, standard library HTTP/JSON, existing Go test suite

---

### Task 1: Harden Config Path And Editor Errors

**Files:**
- Modify: `internal/cfg/cfg.go`
- Modify: `internal/editor/editor.go`
- Modify: `internal/app/app.go`
- Test: `internal/editor/editor_test.go`
- Test: `internal/cfg/cfg_test.go`

- [ ] **Step 1: Write the failing config path tests**

```go
package cfg

import "testing"

func TestConfigDirFromBasePath(t *testing.T) {
	got := configDirFromBase("/tmp/config-root")
	want := "/tmp/config-root/hxtrainer"
	if got != want {
		t.Fatalf("configDirFromBase() = %q, want %q", got, want)
	}
}

func TestFallbackConfigBaseUsesHomeDirectory(t *testing.T) {
	got := fallbackConfigBase("linux", "/home/tester")
	want := "/home/tester/.config"
	if got != want {
		t.Fatalf("fallbackConfigBase() = %q, want %q", got, want)
	}
}
```

- [ ] **Step 2: Run the config tests to verify they fail**

Run: `go test ./internal/cfg -run 'TestConfigDirFromBasePath|TestFallbackConfigBaseUsesHomeDirectory'`
Expected: FAIL with undefined helpers.

- [ ] **Step 3: Write the failing editor tests**

```go
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
```

- [ ] **Step 4: Run the editor tests to verify they fail**

Run: `go test ./internal/editor -run 'TestOpenEditorReturnsMissingFileError|TestOpenEditorReturnsHelpfulErrorWhenBinaryMissing'`
Expected: FAIL because `OpenEditor` returns `bool` and `ErrHelixNotFound` does not exist.

- [ ] **Step 5: Implement minimal config helpers and wire them into cfg**

```go
const appName = "hxtrainer"

func GetConfigDir() string {
	base, err := os.UserConfigDir()
	if err == nil && base != "" {
		return configDirFromBase(base)
	}

	home, _ := os.UserHomeDir()
	return configDirFromBase(fallbackConfigBase(runtime.GOOS, home))
}

func configDirFromBase(base string) string {
	return filepath.Join(base, appName)
}

func fallbackConfigBase(goos, home string) string {
	if home == "" {
		return "." + appName
	}
	switch goos {
	case "windows":
		return home
	default:
		return filepath.Join(home, ".config")
	}
}
```

- [ ] **Step 6: Implement minimal editor error API**

```go
var ErrHelixNotFound = errors.New("helix executable not found")

func OpenEditor(filePath string, cwd string) error {
	if _, err := os.Stat(filePath); err != nil {
		return err
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	path, err := exec.LookPath("hx")
	if err != nil {
		return fmt.Errorf("%w: install Helix and make sure hx is in PATH", ErrHelixNotFound)
	}

	cmd := exec.Command(path, absPath)
	if cwd != "" {
		cmd.Dir = cwd
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run hx: %w", err)
	}
	return nil
}
```

- [ ] **Step 7: Update app call sites to consume `error` instead of `bool`**

```go
if err := editor.OpenEditor(tmpFilePath, tmpDir); err != nil {
	fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Ошибка запуска редактора Helix: %v", err))
	os.RemoveAll(tmpDir)
	return
}
```

- [ ] **Step 8: Run the targeted tests to verify they pass**

Run: `go test ./internal/cfg ./internal/editor -run 'TestConfigDirFromBasePath|TestFallbackConfigBaseUsesHomeDirectory|TestOpenEditorReturnsMissingFileError|TestOpenEditorReturnsHelpfulErrorWhenBinaryMissing'`
Expected: PASS

- [ ] **Step 9: Commit**

```bash
git add internal/cfg/cfg.go internal/cfg/cfg_test.go internal/editor/editor.go internal/editor/editor_test.go internal/app/app.go
git commit -m "fix: harden config and editor runtime errors"
```

### Task 2: Replace Recursive Play Flow With Testable Action Loop

**Files:**
- Modify: `internal/app/app.go`
- Test: `internal/app/app_test.go`

- [ ] **Step 1: Write the failing play-action tests**

```go
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
```

- [ ] **Step 2: Run the play-action tests to verify they fail**

Run: `go test ./internal/app -run TestParsePostChallengeAction`
Expected: FAIL with undefined `postChallengeAction` and `parsePostChallengeAction`.

- [ ] **Step 3: Add a helper for parsing actions and prompting line-based input**

```go
type postChallengeAction int

const (
	actionNext postChallengeAction = iota + 1
	actionRetry
	actionQuit
)

func parsePostChallengeAction(input string) (postChallengeAction, error) {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "j":
		return actionNext, nil
	case "k":
		return actionRetry, nil
	case "q":
		return actionQuit, nil
	default:
		return 0, fmt.Errorf("неизвестное действие %q", strings.TrimSpace(input))
	}
}

func readPostChallengeAction(in io.Reader, out io.Writer) postChallengeAction {
	reader := bufio.NewReader(in)
	for {
		fmt.Fprint(out, "Введите действие [j=next, k=retry, q=quit] и нажмите Enter: ")
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return actionQuit
		}
		action, parseErr := parsePostChallengeAction(line)
		if parseErr == nil {
			return action
		}
		fmt.Fprintln(out, "Неизвестная команда. Используйте j, k или q.")
		if err == io.EOF {
			return actionQuit
		}
	}
}
```

- [ ] **Step 4: Replace recursive `runPlay` transition with iterative state updates**

```go
var requestedID string
if len(args) > 0 {
	requestedID = args[0]
}

for {
	selected, err := selectChallenge(challengeList, requestedID)
	if err != nil {
		fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Ошибка выбора челленджа: %v", err))
		return
	}

	retrySame := false
	for {
		// existing workspace/editor/judge flow
		action := readPostChallengeAction(cmd.InOrStdin(), cmd.OutOrStdout())
		if action == actionRetry {
			retrySame = true
			continue
		}
		if action == actionQuit {
			fmt.Fprintln(cmd.OutOrStdout(), "Выход.")
			return
		}
		requestedID = ""
		break
	}

	if retrySame {
		continue
	}
}
```

- [ ] **Step 5: Handle ignored errors in the play path**

```go
if err := copyFile(extraPath, destPath); err != nil {
	fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Ошибка копирования файла %s: %v", filepath.Base(extraPath), err))
	os.RemoveAll(tmpDir)
	return
}

if _, err := database.LogAttempt(selected.ID, isCorrect, duration); err != nil {
	fmt.Fprintln(cmd.ErrOrStderr(), text.FgYellow.Sprintf("Предупреждение: не удалось сохранить попытку: %v", err))
}
```

- [ ] **Step 6: Run the targeted app tests**

Run: `go test ./internal/app -run 'TestParsePostChallengeAction|TestStatsReset'`
Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add internal/app/app.go internal/app/app_test.go
git commit -m "fix: remove recursive play flow"
```

### Task 3: Replace Broken Upgrade Check And Update README

**Files:**
- Modify: `internal/app/app.go`
- Modify: `README.md`
- Test: `internal/app/app_test.go`

- [ ] **Step 1: Write the failing upgrade tests**

```go
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
```

- [ ] **Step 2: Run the upgrade tests to verify they fail**

Run: `go test ./internal/app -run 'TestCompareReleaseVersions|TestLatestVersionFromResponse'`
Expected: FAIL with undefined helpers.

- [ ] **Step 3: Implement minimal GitHub release parsing and semantic comparison**

```go
type githubRelease struct {
	TagName string `json:"tag_name"`
}

func latestVersionFromResponse(r io.Reader) (string, error) {
	var release githubRelease
	if err := json.NewDecoder(r).Decode(&release); err != nil {
		return "", err
	}
	version := strings.TrimPrefix(strings.TrimSpace(release.TagName), "v")
	if version == "" {
		return "", errors.New("empty tag_name in release response")
	}
	return version, nil
}

func compareReleaseVersions(a, b string) int {
	parse := func(v string) []int {
		parts := strings.Split(strings.TrimPrefix(v, "v"), ".")
		out := make([]int, 3)
		for i := 0; i < len(parts) && i < 3; i++ {
			out[i], _ = strconv.Atoi(parts[i])
		}
		return out
	}

	av := parse(a)
	bv := parse(b)
	for i := 0; i < 3; i++ {
		if av[i] > bv[i] {
			return 1
		}
		if av[i] < bv[i] {
			return -1
		}
	}
	return 0
}
```

- [ ] **Step 4: Replace broken `go.mod` scraping with GitHub Releases API**

```go
func getLatestVersion() string {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/70l571y/HelixTrainer/releases/latest", nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	version, err := latestVersionFromResponse(resp.Body)
	if err != nil {
		return ""
	}
	return version
}
```

- [ ] **Step 5: Update command messaging and README copy**

```md
* **Linux**: `~/.config/hxtrainer/` (или `$XDG_CONFIG_HOME/hxtrainer`)
* **macOS**: `~/Library/Application Support/hxtrainer/`
* **Windows**: `%APPDATA%\\hxtrainer\\`
```

```go
if latest == "" {
	fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Не удалось проверить обновления через GitHub Releases."))
	return
}
```

- [ ] **Step 6: Run the targeted upgrade tests**

Run: `go test ./internal/app -run 'TestCompareReleaseVersions|TestLatestVersionFromResponse'`
Expected: PASS

- [ ] **Step 7: Run the full project test suite**

Run: `go test ./...`
Expected: PASS

- [ ] **Step 8: Commit**

```bash
git add internal/app/app.go internal/app/app_test.go README.md
git commit -m "fix: replace broken upgrade version check"
```
