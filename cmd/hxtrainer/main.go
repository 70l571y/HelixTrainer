// Package main - точка входа приложения HelixTrainer.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	challengesdata "github.com/70l571y/HelixTrainer/challenges_data"
	"github.com/70l571y/HelixTrainer/internal/app"
	"github.com/70l571y/HelixTrainer/internal/buildinfo"
	"github.com/70l571y/HelixTrainer/internal/cfg"

	"github.com/spf13/cobra"
)

var (
	verbose bool
)

func newRootCommand(challengesDir string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "hxtrainer",
		Short: "HelixTrainer - Тренажёр для редактора Helix",
		Long: `HelixTrainer - интерактивный тренажёр для освоения редактора кода Helix.
Решайте челленджи, отрабатывая навыки работы с Helix.`,
		Example: `  hxtrainer play
  hxtrainer play --track core --strategy weak-skills
  hxtrainer list --tag movement_basic
  hxtrainer list --json
  hxtrainer stats
  hxtrainer stats --difficulty medium
  hxtrainer stats export attempts.json
  hxtrainer history hello_world --json
  hxtrainer queue --strategy weak-skills
  hxtrainer stats reset
  hxtrainer stats reset --yes
  hxtrainer doctor
  hxtrainer completion bash`,
		Version: buildinfo.CurrentVersion(),
	}

	app.InitCommands(rootCmd, challengesDir)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Подробный вывод")
	return rootCmd
}

func main() {

	execPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка получения пути: %v\n", err)
		os.Exit(1)
	}

	execDir := filepath.Dir(execPath)
	challengesDir := filepath.Join(execDir, "challenges_data", "go")

	if _, err := os.Stat(challengesDir); os.IsNotExist(err) {
		cwd, _ := os.Getwd()
		challengesDir = filepath.Join(cwd, "challenges_data", "go")
	}

	if _, err := os.Stat(challengesDir); os.IsNotExist(err) {
		configChallengesRoot := cfg.GetChallengesRootDir()
		if err := challengesdata.SyncToDir(configChallengesRoot); err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка подготовки встроенных челленджей: %v\n", err)
			os.Exit(1)
		}
		challengesDir = cfg.GetChallengesDir()
	}

	rootCmd := newRootCommand(challengesDir)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
