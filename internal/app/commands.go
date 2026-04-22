package app

import (
	"net/http"
	"os/exec"

	"github.com/70l571y/HelixTrainer/internal/database"
	"github.com/spf13/cobra"
)

var (
	challengesDir     string
	verbose           bool
	resetAttempts     = database.ResetAttempts
	releaseHTTPClient = http.DefaultClient
	binaryLookPath    = exec.LookPath
)

const githubLatestReleaseURL = "https://api.github.com/repos/70l571y/HelixTrainer/releases/latest"

type postChallengeAction int

const (
	actionNext postChallengeAction = iota + 1
	actionRetry
	actionQuit
)

func InitCommands(rootCmd *cobra.Command, cd string) {
	challengesDir = cd

	rootCmd.AddCommand(playCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(queueCmd)
	rootCmd.AddCommand(historyCmd)
	rootCmd.AddCommand(newCompletionCmd(rootCmd))
}

func init() {
	statsCmd.AddCommand(statsResetCmd)
	statsCmd.AddCommand(statsExportCmd)
	statsCmd.AddCommand(statsImportCmd)
	_ = statsResetCmd.Flags().Bool("yes", false, "Подтвердить сброс без prompt")
	_ = listCmd.Flags().Bool("json", false, "Вывести список в JSON")
	_ = statsCmd.Flags().Bool("json", false, "Вывести статистику в JSON")
	_ = playCmd.Flags().String("difficulty", "", "Фильтр по сложности: easy|medium|hard")
	_ = listCmd.Flags().String("difficulty", "", "Фильтр по сложности: easy|medium|hard")
	_ = statsCmd.Flags().String("difficulty", "", "Фильтр по сложности: easy|medium|hard")
	_ = playCmd.Flags().StringSlice("tag", nil, "Фильтр по тегам (можно указать несколько)")
	_ = listCmd.Flags().StringSlice("tag", nil, "Фильтр по тегам (можно указать несколько)")
	_ = statsCmd.Flags().StringSlice("tag", nil, "Фильтр по тегам (можно указать несколько)")
	_ = playCmd.Flags().String("track", "", "Фильтр по треку: core|optional")
	_ = listCmd.Flags().String("track", "", "Фильтр по треку: core|optional")
	_ = statsCmd.Flags().String("track", "", "Фильтр по треку: core|optional")
	_ = playCmd.Flags().String("strategy", "smart", "Стратегия выбора: smart|progression|weak-skills")
	_ = queueCmd.Flags().Bool("json", false, "Вывести очередь в JSON")
	_ = historyCmd.Flags().Bool("json", false, "Вывести историю в JSON")
	_ = doctorCmd.Flags().Bool("json", false, "Вывести диагностику в JSON")
	_ = queueCmd.Flags().Int("limit", 5, "Сколько челленджей показать")
	_ = historyCmd.Flags().Int("limit", 20, "Сколько попыток показать")
	_ = queueCmd.Flags().String("difficulty", "", "Фильтр по сложности: easy|medium|hard")
	_ = queueCmd.Flags().StringSlice("tag", nil, "Фильтр по тегам (можно указать несколько)")
	_ = queueCmd.Flags().String("track", "", "Фильтр по треку: core|optional")
	_ = queueCmd.Flags().String("strategy", "smart", "Стратегия выбора: smart|progression|weak-skills")
	_ = historyCmd.Flags().String("difficulty", "", "Фильтр по сложности: easy|medium|hard")
	_ = historyCmd.Flags().StringSlice("tag", nil, "Фильтр по тегам (можно указать несколько)")
	_ = historyCmd.Flags().String("track", "", "Фильтр по треку: core|optional")
	_ = statsImportCmd.Flags().Bool("replace", false, "Очистить текущую статистику перед импортом")
}
