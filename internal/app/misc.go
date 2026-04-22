package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/70l571y/HelixTrainer/internal/buildinfo"
	"github.com/70l571y/HelixTrainer/internal/cfg"
	"github.com/70l571y/HelixTrainer/internal/challenges"
	"github.com/70l571y/HelixTrainer/internal/database"
	"github.com/70l571y/HelixTrainer/internal/judge"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

func runUpgrade(cmd *cobra.Command, args []string) {
	out := cmd.OutOrStdout()
	errOut := cmd.ErrOrStderr()
	fmt.Fprintln(out, text.FgCyan.Sprintf("Проверка обновлений..."))
	latest := getLatestVersion()
	current := getCurrentVersion()
	if latest == "" {
		fmt.Fprintln(errOut, text.FgRed.Sprintf("Не удалось проверить обновления через GitHub Releases."))
		return
	}
	if current != "devel" && compareReleaseVersions(current, latest) >= 0 {
		fmt.Fprintln(out, text.FgGreen.Sprintf("Уже установлена последняя версия (%s).", current))
		return
	}
	fmt.Fprintln(out, text.FgYellow.Sprintf("Доступно обновление: %s -> %s.", current, latest))
	fmt.Fprintln(out, text.FgYellow.Sprintf("Для обновления выполните:"))
	fmt.Fprintln(out, "  go install github.com/70l571y/HelixTrainer/cmd/hxtrainer@latest")
}

func runDoctor(cmd *cobra.Command, args []string) {
	checks := collectDoctorChecks(cfg.GetConfigDir(), cfg.GetDBPath(), challengesDir)
	jsonOutput, _ := cmd.Flags().GetBool("json")
	if jsonOutput {
		writeJSON(cmd.OutOrStdout(), checks)
		return
	}
	t := table.NewWriter()
	t.SetOutputMirror(cmd.OutOrStdout())
	t.AppendHeader(table.Row{"Проверка", "Статус", "Детали"})
	for _, check := range checks {
		status := text.FgGreen.Sprintf("OK")
		if !check.OK {
			status = text.FgRed.Sprintf("FAIL")
		}
		t.AppendRow(table.Row{check.Name, status, check.Details})
	}
	t.Render()
}

func displayFeedback(userText, goalText, judgeMode string) {
	if judge.CheckSolution(userText, goalText, judgeMode) {
		fmt.Println("🎉 Успех! Решение верное.")
		return
	}
	fmt.Println("❌ Решение неверное. Вот diff:")
	fmt.Println(judge.GenerateDiff(userText, goalText))
	if judgeMode == "ast" || judgeMode == "go_ast" {
		fmt.Println(text.FgYellow.Sprintf("Примечание: Режим AST. Структура должна совпадать точно, но форматирование может отличаться."))
	}
}

func challengeStatusLabel(completed bool) string {
	if completed {
		return "completed"
	}
	return "unsolved"
}

func challengeFilterFromCommand(cmd *cobra.Command) challenges.ChallengeFilter {
	difficulty, _ := cmd.Flags().GetString("difficulty")
	tags, _ := cmd.Flags().GetStringSlice("tag")
	track, _ := cmd.Flags().GetString("track")
	return challenges.ChallengeFilter{Difficulty: difficulty, Tags: tags, Track: track}
}

func hasChallengeFilter(cmd *cobra.Command) bool {
	filter := challengeFilterFromCommand(cmd)
	return strings.TrimSpace(filter.Difficulty) != "" || strings.TrimSpace(filter.Track) != "" || len(filter.Tags) > 0
}

func writeJSON(w io.Writer, payload any) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(payload); err != nil {
		fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Ошибка JSON вывода: %v", err))
	}
}

func collectDoctorChecks(configDir, dbPath, challengesPath string) []doctorCheck {
	return []doctorCheck{
		binaryCheck("helix", "hx"),
		binaryCheck("git", "git"),
		pathCheck("config_dir", configDir),
		pathCheck("db_parent", filepath.Dir(dbPath)),
		pathCheck("challenges_dir", challengesPath),
	}
}

func binaryCheck(name, binary string) doctorCheck {
	path, err := binaryLookPath(binary)
	if err != nil {
		return doctorCheck{Name: name, OK: false, Details: fmt.Sprintf("%s not found in PATH", binary)}
	}
	return doctorCheck{Name: name, OK: true, Details: path}
}

func pathCheck(name, path string) doctorCheck {
	if path == "" {
		return doctorCheck{Name: name, OK: false, Details: "empty path"}
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		return doctorCheck{Name: name, OK: false, Details: err.Error()}
	}
	return doctorCheck{Name: name, OK: true, Details: path}
}

func buildQueue(challengeList []challenges.Challenge, attempts []database.Attempt, strategy string, limit int) []queueJSONEntry {
	if limit <= 0 {
		limit = 5
	}
	if len(challengeList) == 0 {
		return []queueJSONEntry{}
	}
	ordered := orderedChallengesForStrategy(challengeList, attempts, strategy)
	if len(ordered) > limit {
		ordered = ordered[:limit]
	}
	reason := queueReason(strategy)
	queue := make([]queueJSONEntry, 0, len(ordered))
	for _, challenge := range ordered {
		queue = append(queue, queueJSONEntry{ID: challenge.ID, Difficulty: challenge.Difficulty, Tags: challenge.Tags, Reason: reason})
	}
	return queue
}

func orderedChallengesForStrategy(challengeList []challenges.Challenge, attempts []database.Attempt, strategy string) []challenges.Challenge {
	switch stringsLowerTrim(strategy) {
	case "progression":
		return sortByProgressionWithSolvedLast(challengeList, attempts)
	case "weak-skills", "weak_skills", "weak":
		return sortByWeakSkills(challengeList, attempts)
	default:
		return sortByProgressionWithSolvedLast(challengeList, attempts)
	}
}

func queueReason(strategy string) string {
	switch stringsLowerTrim(strategy) {
	case "progression":
		return "progression"
	case "weak-skills", "weak_skills", "weak":
		return "weak_skills"
	default:
		return "smart"
	}
}

func sortByProgressionWithSolvedLast(challengeList []challenges.Challenge, attempts []database.Attempt) []challenges.Challenge {
	solved := map[string]bool{}
	for _, attempt := range attempts {
		if attempt.IsCorrect {
			solved[attempt.ChallengeID] = true
		}
	}
	var unsolved, done []challenges.Challenge
	for _, challenge := range challengeList {
		if solved[challenge.ID] {
			done = append(done, challenge)
		} else {
			unsolved = append(unsolved, challenge)
		}
	}
	challenges.SortForProgression(unsolved)
	challenges.SortForProgression(done)
	return append(unsolved, done...)
}

func sortByWeakSkills(challengeList []challenges.Challenge, attempts []database.Attempt) []challenges.Challenge {
	remaining := append([]challenges.Challenge(nil), challengeList...)
	var ordered []challenges.Challenge
	for len(remaining) > 0 {
		next := challenges.SelectWeakestChallenge(remaining, attempts)
		if next.ID == "" {
			break
		}
		ordered = append(ordered, next)
		filtered := remaining[:0]
		for _, challenge := range remaining {
			if challenge.ID != next.ID {
				filtered = append(filtered, challenge)
			}
		}
		remaining = filtered
	}
	return ordered
}

func getCurrentVersion() string {
	return buildinfo.CurrentVersion()
}

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func getLatestVersion() string {
	req, err := http.NewRequest(http.MethodGet, githubLatestReleaseURL, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := releaseHTTPClient.Do(req)
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
		parts := strings.Split(strings.TrimPrefix(strings.TrimSpace(v), "v"), ".")
		out := make([]int, 3)
		for i := 0; i < len(parts) && i < 3; i++ {
			n, err := strconv.Atoi(parts[i])
			if err == nil {
				out[i] = n
			}
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

func newCompletionCmd(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:       "completion [bash|zsh|fish|powershell]",
		Short:     "Сгенерировать shell completion",
		Args:      cobra.ExactValidArgs(1),
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return rootCmd.GenBashCompletion(cmd.OutOrStdout())
			case "zsh":
				return rootCmd.GenZshCompletion(cmd.OutOrStdout())
			case "fish":
				return rootCmd.GenFishCompletion(cmd.OutOrStdout(), true)
			case "powershell":
				return rootCmd.GenPowerShellCompletionWithDesc(cmd.OutOrStdout())
			default:
				return fmt.Errorf("unsupported shell %q", args[0])
			}
		},
	}
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, sourceFile)
	return err
}

func stringsLowerTrim(s string) string  { return strings.ToLower(strings.TrimSpace(s)) }
func stringsJoin(items []string) string { return strings.Join(items, ", ") }
