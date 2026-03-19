// Package main - точка входа приложения HelixTrainer.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	challengesdata "github.com/70l571y/HelixTrainer/challenges_data"
	"github.com/70l571y/HelixTrainer/internal/app"
	"github.com/70l571y/HelixTrainer/internal/cfg"

	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
	verbose bool
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "hxtrainer",
		Short: "HelixTrainer - Тренажёр для редактора Helix",
		Long: `HelixTrainer - интерактивный тренажёр для освоения редактора кода Helix.
Решайте челленджи, отрабатывая навыки работы с Helix.`,
		Version: version,
	}

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

	app.InitCommands(rootCmd, challengesDir)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Подробный вывод")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
