package challenges

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestChallengeDataIsConsistent(t *testing.T) {
	root := filepath.Join("..", "..", "challenges_data", "go")
	labelSet := helixLabelSet(t)

	challengeList, err := LoadChallenges(root)
	if err != nil {
		t.Fatalf("LoadChallenges(%q) error = %v", root, err)
	}
	if len(challengeList) == 0 {
		t.Fatalf("LoadChallenges(%q) returned no challenges", root)
	}

	for _, challenge := range challengeList {
		challenge := challenge
		t.Run(challenge.ID, func(t *testing.T) {
			for _, tag := range challenge.Tags {
				if _, ok := labelSet[tag]; !ok {
					t.Fatalf("%s uses tag %q, but docs/HelixLabels.md does not define it", challenge.ID, tag)
				}
			}

			assertFileExists(t, challenge.StartPath)
			for _, goalPath := range challenge.GoalPaths {
				assertFileExists(t, goalPath)
			}
			for _, extraPath := range challenge.ExtraFilePaths {
				assertFileExists(t, extraPath)
			}
			for _, dirtyFixture := range challenge.GitDirtyFiles {
				assertFileExists(t, filepath.Join(challenge.DirPath, dirtyFixture))
			}
			for _, goalPath := range challenge.ValidationMap {
				assertFileExists(t, goalPath)
			}

			goFiles, err := filepath.Glob(filepath.Join(challenge.DirPath, "*.go"))
			if err != nil {
				t.Fatalf("Glob(%q) error = %v", challenge.DirPath, err)
			}
			for _, path := range goFiles {
				content, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("ReadFile(%q) error = %v", path, err)
				}
				if !strings.HasPrefix(string(content), "//go:build ignore\n") {
					t.Fatalf("%s must start with //go:build ignore to stay out of normal package builds", path)
				}
			}

			if !strings.Contains(strings.Join(challenge.Tags, ","), "register_blackhole") {
				return
			}
			if strings.Contains(challenge.Tips, "'_d'") {
				t.Fatalf("%s tips mention '_d', but modern Helix requires register selection via '\"_d' or a documented shortcut", challenge.ID)
			}
			if !strings.Contains(challenge.Tips, "'\"_d'") && !strings.Contains(challenge.Tips, "'Alt-d'") {
				t.Fatalf("%s tips must mention either '\"_d' or 'Alt-d' for blackhole delete", challenge.ID)
			}
		})
	}
}

func helixLabelSet(t *testing.T) map[string]struct{} {
	t.Helper()

	path := filepath.Join("..", "..", "docs", "HelixLabels.md")
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("Open(%q) error = %v", path, err)
	}
	defer file.Close()

	labelPattern := regexp.MustCompile("`([^`]+)`")
	labels := make(map[string]struct{})

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := labelPattern.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			labels[match[1]] = struct{}{}
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("Scan(%q) error = %v", path, err)
	}

	if len(labels) == 0 {
		t.Fatalf("no labels parsed from %q", path)
	}

	return labels
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file %q to exist: %v", path, err)
	}
}
