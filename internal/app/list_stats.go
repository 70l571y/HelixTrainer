package app

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/70l571y/HelixTrainer/internal/challenges"
	"github.com/70l571y/HelixTrainer/internal/database"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

func runList(cmd *cobra.Command, args []string) {
	jsonOutput, _ := cmd.Flags().GetBool("json")

	if err := database.InitDB(); err != nil {
		fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Ошибка инициализации БД: %v", err))
		return
	}

	attempts, _ := database.GetAllAttempts()
	challengeList, err := challenges.LoadChallenges(challengesDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Ошибка загрузки челленджей: %v", err))
		return
	}
	challengeList = challenges.FilterChallenges(challengeList, challengeFilterFromCommand(cmd))
	challenges.SortForProgression(challengeList)

	completedIDs := make(map[string]bool)
	for _, a := range attempts {
		if a.IsCorrect {
			completedIDs[a.ChallengeID] = true
		}
	}

	if jsonOutput {
		payload := make([]listJSONEntry, 0, len(challengeList))
		for _, c := range challengeList {
			payload = append(payload, listJSONEntry{
				ID:         c.ID,
				Difficulty: c.Difficulty,
				Language:   c.Language,
				Tags:       c.Tags,
				Status:     challengeStatusLabel(completedIDs[c.ID]),
				Completed:  completedIDs[c.ID],
			})
		}
		writeJSON(cmd.OutOrStdout(), payload)
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Сложность", "Язык", "Метки", "Статус"})
	for _, c := range challengeList {
		status := "Не решено"
		displayID := c.ID
		if completedIDs[c.ID] {
			displayID = text.FgGreen.Sprintf("%s", c.ID)
			status = text.Bold.Sprintf("%s", text.FgGreen.Sprintf("Выполнено"))
		}
		t.AppendRow(table.Row{displayID, c.Difficulty, c.Language, stringsJoin(c.Tags), status})
	}
	t.Render()
}

func runStats(cmd *cobra.Command, args []string) {
	out := cmd.OutOrStdout()
	jsonOutput, _ := cmd.Flags().GetBool("json")

	if err := database.InitDB(); err != nil {
		fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Ошибка инициализации БД: %v", err))
		return
	}

	attempts, _ := database.GetAllAttempts()
	challengeList, err := challenges.LoadChallenges(challengesDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Ошибка загрузки челленджей: %v", err))
		return
	}
	challengeList = challenges.FilterChallenges(challengeList, challengeFilterFromCommand(cmd))
	allowedChallengeIDs := make(map[string]bool, len(challengeList))
	for _, challenge := range challengeList {
		allowedChallengeIDs[challenge.ID] = true
	}
	if len(challengeList) == 0 {
		if !jsonOutput && !hasChallengeFilter(cmd) {
			fmt.Fprintln(out, text.FgYellow.Sprintf("Попыток пока нет."))
			return
		}
		if jsonOutput {
			writeJSON(out, statsJSONPayload{RecentAttempts: []statsRecentAttemptJSON{}, Challenges: []statsChallengeJSON{}})
			return
		}
		fmt.Fprintln(out, text.FgYellow.Sprintf("Под выбранные фильтры челленджи не найдены."))
		return
	}

	filteredAttempts := make([]database.Attempt, 0, len(attempts))
	for _, attempt := range attempts {
		if allowedChallengeIDs[attempt.ChallengeID] {
			filteredAttempts = append(filteredAttempts, attempt)
		}
	}
	attempts = filteredAttempts

	if len(attempts) == 0 && !jsonOutput {
		fmt.Fprintln(out, text.FgYellow.Sprintf("Попыток пока нет."))
		return
	}

	sort.Slice(attempts, func(i, j int) bool { return attempts[i].Timestamp.After(attempts[j].Timestamp) })
	attemptsByID := make(map[string][]database.Attempt)
	for _, a := range attempts {
		attemptsByID[a.ChallengeID] = append(attemptsByID[a.ChallengeID], a)
	}
	challenges.SortForProgression(challengeList)

	if jsonOutput {
		payload := statsJSONPayload{
			RecentAttempts: make([]statsRecentAttemptJSON, 0),
			Challenges:     make([]statsChallengeJSON, 0, len(challengeList)),
		}
		for i := 0; i < 20 && i < len(attempts); i++ {
			a := attempts[i]
			result := "pass"
			if !a.IsCorrect {
				result = "fail"
			}
			payload.RecentAttempts = append(payload.RecentAttempts, statsRecentAttemptJSON{
				Timestamp:       a.Timestamp.Format(time.RFC3339),
				ChallengeID:     a.ChallengeID,
				Result:          result,
				DurationSeconds: a.Duration,
			})
		}
		for _, c := range challengeList {
			cAttempts := attemptsByID[c.ID]
			successful := successfulAttempts(cAttempts)
			item := statsChallengeJSON{ID: c.ID, Status: "no_attempts", Attempts: len(cAttempts)}
			if len(cAttempts) > 0 {
				item.Status = "unsolved"
			}
			if len(successful) > 0 {
				item.Status = "completed"
				bestTime := bestAttemptDuration(successful)
				item.BestTimeSeconds = &bestTime
				if ms := challenges.GetMilestone(bestTime, c.AuthorTime); ms.Name != "" {
					item.Milestone = ms.Name
				}
			}
			payload.Challenges = append(payload.Challenges, item)
		}
		writeJSON(out, payload)
		return
	}

	fmt.Fprintln(out, text.Bold.Sprintf("Последняя активность (20)"))
	recentTable := table.NewWriter()
	recentTable.SetOutputMirror(out)
	recentTable.AppendHeader(table.Row{"Время", "Челлендж", "Результат", "Время"})
	for i := 0; i < 20 && i < len(attempts); i++ {
		a := attempts[i]
		result := text.FgGreen.Sprintf("Pass")
		if !a.IsCorrect {
			result = text.FgRed.Sprintf("Fail")
		}
		recentTable.AppendRow(table.Row{a.Timestamp.Format("2006-01-02 15:04"), a.ChallengeID, result, fmt.Sprintf("%.2fс", a.Duration)})
	}
	recentTable.Render()
	fmt.Fprintln(out)

	fmt.Fprintln(out, text.Bold.Sprintf("Прогресс челленджей"))
	progressTable := table.NewWriter()
	progressTable.SetOutputMirror(out)
	progressTable.AppendHeader(table.Row{"Челлендж", "Статус", "Лучшее время", "Веха", "Попыток"})
	for _, c := range challengeList {
		cAttempts := attemptsByID[c.ID]
		successful := successfulAttempts(cAttempts)
		status := "Нет попыток"
		if len(cAttempts) > 0 {
			status = "Не решено"
		}
		bestTimeStr := "-"
		msDisplay := "-"
		if len(successful) > 0 {
			status = text.Bold.Sprintf("%s", text.FgGreen.Sprintf("Выполнено"))
			bestTime := bestAttemptDuration(successful)
			bestTimeStr = text.FgGreen.Sprintf("%.2fс", bestTime)
			if ms := challenges.GetMilestone(bestTime, c.AuthorTime); ms.Name != "" {
				msDisplay = fmt.Sprintf("%s %s", ms.Symbol, ms.Name)
			}
		}
		progressTable.AppendRow(table.Row{c.ID, status, bestTimeStr, msDisplay, len(cAttempts)})
	}
	progressTable.Render()
}

func runStatsReset(cmd *cobra.Command, args []string) error {
	yes, _ := cmd.Flags().GetBool("yes")
	if !yes {
		confirmed, err := confirmStatsReset(cmd.InOrStdin(), cmd.OutOrStdout())
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Fprintln(cmd.OutOrStdout(), "Сброс статистики отменён.")
			return nil
		}
	}
	removed, err := resetAttempts()
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Удалено %d попыток из статистики.\n", removed)
	return nil
}

func runStatsExport(cmd *cobra.Command, args []string) error {
	if err := database.InitDB(); err != nil {
		return err
	}
	file, err := os.Create(args[0])
	if err != nil {
		return err
	}
	defer file.Close()
	if err := database.ExportAttempts(file); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Статистика экспортирована в %s.\n", args[0])
	return nil
}

func runStatsImport(cmd *cobra.Command, args []string) error {
	if err := database.InitDB(); err != nil {
		return err
	}
	replace, _ := cmd.Flags().GetBool("replace")
	file, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer file.Close()
	count, err := database.ImportAttempts(file, replace)
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Импортировано %d попыток из %s.\n", count, args[0])
	return nil
}

func runQueue(cmd *cobra.Command, args []string) {
	jsonOutput, _ := cmd.Flags().GetBool("json")
	limit, _ := cmd.Flags().GetInt("limit")
	strategy, _ := cmd.Flags().GetString("strategy")
	challengeList, err := challenges.LoadChallenges(challengesDir)
	if err != nil {
		fmt.Fprintln(cmd.ErrOrStderr(), text.FgRed.Sprintf("Ошибка загрузки челленджей: %v", err))
		return
	}
	challengeList = challenges.FilterChallenges(challengeList, challengeFilterFromCommand(cmd))
	attempts, _ := database.GetAllAttempts()
	queue := buildQueue(challengeList, attempts, strategy, limit)
	if jsonOutput {
		writeJSON(cmd.OutOrStdout(), queue)
		return
	}
	t := table.NewWriter()
	t.SetOutputMirror(cmd.OutOrStdout())
	t.AppendHeader(table.Row{"ID", "Сложность", "Причина", "Метки"})
	for _, item := range queue {
		t.AppendRow(table.Row{item.ID, item.Difficulty, item.Reason, stringsJoin(item.Tags)})
	}
	t.Render()
}

func runHistory(cmd *cobra.Command, args []string) {
	jsonOutput, _ := cmd.Flags().GetBool("json")
	limit, _ := cmd.Flags().GetInt("limit")
	var challengeID string
	if len(args) > 0 {
		challengeID = args[0]
	}

	attempts, err := database.GetAttempts(challengeID)
	if err != nil {
		fmt.Fprintln(cmd.ErrOrStderr(), text.FgRed.Sprintf("Ошибка чтения истории: %v", err))
		return
	}
	if challengeID == "" {
		challengeList, err := challenges.LoadChallenges(challengesDir)
		if err == nil {
			filtered := challenges.FilterChallenges(challengeList, challengeFilterFromCommand(cmd))
			allowed := map[string]bool{}
			for _, challenge := range filtered {
				allowed[challenge.ID] = true
			}
			if len(filtered) > 0 {
				kept := attempts[:0]
				for _, attempt := range attempts {
					if allowed[attempt.ChallengeID] {
						kept = append(kept, attempt)
					}
				}
				attempts = kept
			}
		}
	}
	if limit > 0 && len(attempts) > limit {
		attempts = attempts[:limit]
	}
	if jsonOutput {
		payload := make([]historyJSONEntry, 0, len(attempts))
		for _, attempt := range attempts {
			payload = append(payload, historyJSONEntry{
				ChallengeID:     attempt.ChallengeID,
				Timestamp:       attempt.Timestamp.Format(time.RFC3339),
				IsCorrect:       attempt.IsCorrect,
				DurationSeconds: attempt.Duration,
			})
		}
		writeJSON(cmd.OutOrStdout(), payload)
		return
	}
	t := table.NewWriter()
	t.SetOutputMirror(cmd.OutOrStdout())
	t.AppendHeader(table.Row{"Время", "Челлендж", "Результат", "Длительность"})
	for _, attempt := range attempts {
		result := text.FgGreen.Sprintf("Pass")
		if !attempt.IsCorrect {
			result = text.FgRed.Sprintf("Fail")
		}
		t.AppendRow(table.Row{attempt.Timestamp.Format("2006-01-02 15:04"), attempt.ChallengeID, result, fmt.Sprintf("%.2fс", attempt.Duration)})
	}
	t.Render()
}

func confirmStatsReset(in io.Reader, out io.Writer) (bool, error) {
	fmt.Fprint(out, "Внимание: будут удалены все попытки и рекорды из статистики. Сбросить? [y/N]: ")
	reader := bufio.NewReader(in)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return false, err
	}
	answer := stringsLowerTrim(line)
	if err == io.EOF && answer == "" {
		return false, errors.New("подтверждение для сброса статистики не получено")
	}
	return answer == "y" || answer == "yes", nil
}

func successfulAttempts(attempts []database.Attempt) []database.Attempt {
	var successful []database.Attempt
	for _, attempt := range attempts {
		if attempt.IsCorrect {
			successful = append(successful, attempt)
		}
	}
	return successful
}

func bestAttemptDuration(attempts []database.Attempt) float64 {
	best := attempts[0].Duration
	for _, attempt := range attempts[1:] {
		if attempt.Duration < best {
			best = attempt.Duration
		}
	}
	return best
}
