package app

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/70l571y/HelixTrainer/internal/challenges"
	"github.com/70l571y/HelixTrainer/internal/database"
	"github.com/70l571y/HelixTrainer/internal/editor"
	"github.com/70l571y/HelixTrainer/internal/judge"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

func runPlay(cmd *cobra.Command, args []string) {
	if err := database.InitDB(); err != nil {
		fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Ошибка инициализации БД: %v", err))
		return
	}

	challengeList, err := challenges.LoadChallenges(challengesDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Ошибка загрузки челленджей: %v", err))
		return
	}

	if len(challengeList) == 0 {
		fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Челленджи не найдены!"))
		return
	}

	filter := challengeFilterFromCommand(cmd)
	challengeList = challenges.FilterChallenges(challengeList, filter)
	if len(challengeList) == 0 {
		fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Нет челленджей, подходящих под выбранные фильтры."))
		return
	}

	requestedID := ""
	if len(args) > 0 {
		requestedID = strings.TrimSpace(args[0])
	}
	strategy, _ := cmd.Flags().GetString("strategy")

	for {
		allAttempts, _ := database.GetAllAttempts()
		selected, err := selectChallengeForStrategy(challengeList, requestedID, strategy, allAttempts)
		if err != nil {
			fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Ошибка выбора челленджа: %v", err))
			return
		}

		fmt.Println(text.Bold.Sprintf("Запуск челленджа: %s", selected.ID))
		fmt.Println(selected.Description)

		if _, err := os.Stat(selected.StartPath); os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Файл start не найден: %s", selected.StartPath))
			return
		}

		if len(selected.GoalPaths) == 0 {
			fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Файлы goal не найдены для челленджа %s", selected.ID))
			return
		}

		startContent, err := os.ReadFile(selected.StartPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Ошибка чтения start файла: %v", err))
			return
		}

		var goalContents []string
		for _, gp := range selected.GoalPaths {
			content, err := os.ReadFile(gp)
			if err == nil {
				goalContents = append(goalContents, challenges.BuildFileContent(selected, string(content)))
			}
		}

		if len(goalContents) == 0 {
			fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Файлы goal существуют, но не удалось прочитать!"))
			return
		}

		formattedStart := challenges.BuildFileContent(selected, string(startContent))

		for {
			tmpDir, err := os.MkdirTemp("", "helix_trainer*")
			if err != nil {
				fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Ошибка создания временной директории: %v", err))
				return
			}

			mainFileName := challengeMainFileName(selected)
			tmpFilePath := filepath.Join(tmpDir, mainFileName)

			if err := os.WriteFile(tmpFilePath, []byte(formattedStart), 0644); err != nil {
				fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Ошибка записи временного файла: %v", err))
				os.RemoveAll(tmpDir)
				return
			}

			for _, extraPath := range selected.ExtraFilePaths {
				if _, err := os.Stat(extraPath); err == nil {
					destPath := filepath.Join(tmpDir, filepath.Base(extraPath))
					if err := copyFile(extraPath, destPath); err != nil {
						fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Ошибка копирования файла %s: %v", filepath.Base(extraPath), err))
						os.RemoveAll(tmpDir)
						return
					}
				}
			}

			if err := prepareGitWorkspace(selected, tmpDir, mainFileName); err != nil {
				fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Ошибка подготовки git workspace: %v", err))
				os.RemoveAll(tmpDir)
				return
			}

			startTime := time.Now()
			if err := editor.OpenEditor(tmpFilePath, tmpDir); err != nil {
				fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Ошибка запуска редактора Helix: %v", err))
				os.RemoveAll(tmpDir)
				return
			}
			duration := time.Since(startTime).Seconds()

			userContent, err := os.ReadFile(tmpFilePath)
			if err != nil {
				userContent = []byte{}
			}

			judgeMode := selected.JudgeMode
			if judgeMode == "" {
				judgeMode = "exact"
			}

			isCorrect := false
			var matchingGoal string
			feedbackContent := string(userContent)

			if len(selected.ValidationMap) > 0 {
				allPassed := true
				var primaryGoal string

				for filename, goalPath := range selected.ValidationMap {
					targetFile := filepath.Join(tmpDir, filename)
					if _, err := os.Stat(targetFile); os.IsNotExist(err) {
						allPassed = false
						break
					}

					userFileContent, _ := os.ReadFile(targetFile)
					goalRaw, err := os.ReadFile(goalPath)
					if err != nil {
						allPassed = false
						break
					}

					var goalContent string
					if filename == mainFileName {
						goalContent = challenges.BuildFileContent(selected, string(goalRaw))
						primaryGoal = goalContent
					} else {
						goalContent = string(goalRaw)
					}

					if !judge.CheckSolution(string(userFileContent), goalContent, judgeMode) {
						allPassed = false
						feedbackContent = string(userFileContent)
						matchingGoal = goalContent
						break
					}
				}

				if allPassed {
					isCorrect = true
					if primaryGoal != "" {
						matchingGoal = primaryGoal
					} else if len(goalContents) > 0 {
						matchingGoal = goalContents[0]
					}
				}
			} else {
				for _, gContent := range goalContents {
					if judge.CheckSolution(string(userContent), gContent, judgeMode) {
						isCorrect = true
						matchingGoal = gContent
						break
					}
				}
			}

			prevAttempts, _ := database.GetAttempts(selected.ID)
			var prevSuccesses []database.Attempt
			for _, a := range prevAttempts {
				if a.IsCorrect {
					prevSuccesses = append(prevSuccesses, a)
				}
			}

			isNewRecord := false
			prevBestTime := float64(1e9)
			if len(prevSuccesses) > 0 {
				for _, a := range prevSuccesses {
					if a.Duration < prevBestTime {
						prevBestTime = a.Duration
					}
				}
				if duration < prevBestTime {
					isNewRecord = true
				}
			}

			ms := challenges.GetMilestone(duration, selected.AuthorTime)
			prevMsRank := 0
			if len(prevSuccesses) > 0 && selected.AuthorTime > 0 {
				prevMs := challenges.GetMilestone(prevBestTime, selected.AuthorTime)
				prevMsRank = prevMs.Rank
			}
			isNewMilestone := ms.Rank > prevMsRank

			if _, err := database.LogAttempt(selected.ID, isCorrect, duration); err != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), text.FgYellow.Sprintf("Предупреждение: не удалось сохранить попытку: %v", err))
			}

			if isCorrect {
				displayFeedback(feedbackContent, matchingGoal, judgeMode)

				msg := text.Bold.Sprintf("%s", text.FgGreen.Sprintf("Вы выполнили задание за %.2fс", duration))
				if isNewRecord {
					msg += text.Bold.Sprintf("%s", text.FgYellow.Sprintf(" 🏆 Новый рекорд!"))
				}
				if isNewMilestone {
					msg += text.Bold.Sprintf("%s", text.FgMagenta.Sprintf("\n🎉 Новая веха: %s %s!", ms.Symbol, ms.Name))
				} else if ms.Name != "" {
					msg += fmt.Sprintf(" [%s %s]", ms.Symbol, ms.Name)
				}
				fmt.Println(msg)
			} else {
				feedbackGoal := matchingGoal
				if feedbackGoal == "" && len(goalContents) > 0 {
					feedbackGoal = goalContents[0]
				}
				displayFeedback(feedbackContent, feedbackGoal, judgeMode)
				fmt.Println(text.Bold.Sprintf("%s", text.FgRed.Sprintf("❌ Не удалось. Время: %.2fс", duration)))
			}

			os.RemoveAll(tmpDir)

			action := readPostChallengeAction(cmd.InOrStdin(), cmd.OutOrStdout())
			switch action {
			case actionRetry:
				continue
			case actionQuit:
				fmt.Fprintln(cmd.OutOrStdout(), "Выход.")
				return
			case actionNext:
				fmt.Fprintln(cmd.OutOrStdout(), text.FgGreen.Sprintf("Загрузка следующего челленджа..."))
				requestedID = ""
				goto nextChallenge
			}
		}

	nextChallenge:
	}
}

func selectChallengeForStrategy(challengeList []challenges.Challenge, requestedID string, strategy string, attempts []database.Attempt) (challenges.Challenge, error) {
	if requestedID != "" {
		for _, c := range challengeList {
			if c.ID == requestedID {
				return c, nil
			}
		}
		return challenges.Challenge{}, fmt.Errorf("челлендж с ID %s не найден", requestedID)
	}

	switch strings.ToLower(strings.TrimSpace(strategy)) {
	case "", "smart":
		return challenges.SelectSmartChallenge(challengeList), nil
	case "progression":
		return challenges.SelectProgressionChallenge(challengeList, attempts), nil
	case "weak-skills", "weak_skills", "weak":
		return challenges.SelectWeakestChallenge(challengeList, attempts), nil
	default:
		return challenges.Challenge{}, fmt.Errorf("неизвестная стратегия %q", strategy)
	}
}

func challengeMainFileName(challenge challenges.Challenge) string {
	if strings.TrimSpace(challenge.MainFileName) != "" {
		return challenge.MainFileName
	}

	ext := filepath.Ext(challenge.StartPath)
	if ext == "" {
		ext = ".go"
	}

	return "challenge" + ext
}

func prepareGitWorkspace(challenge challenges.Challenge, tmpDir, mainFileName string) error {
	if len(challenge.GitDirtyFiles) == 0 {
		return nil
	}

	if _, err := exec.LookPath("git"); err != nil {
		return nil
	}

	gitCommands := [][]string{
		{"init", "-q"},
		{"config", "user.name", "HelixTrainer"},
		{"config", "user.email", "hxtrainer@example.invalid"},
		{"add", "."},
		{"commit", "-q", "-m", "challenge baseline"},
	}

	for _, args := range gitCommands {
		cmd := exec.Command("git", args...)
		cmd.Dir = tmpDir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
		}
	}

	for targetFile, sourceFixture := range challenge.GitDirtyFiles {
		sourcePath := filepath.Join(challenge.DirPath, sourceFixture)
		content, err := os.ReadFile(sourcePath)
		if err != nil {
			return fmt.Errorf("read dirty fixture %q: %w", sourcePath, err)
		}

		targetPath := filepath.Join(tmpDir, targetFile)
		if targetFile == mainFileName {
			content = []byte(challenges.BuildFileContent(challenge, string(content)))
		}

		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			return fmt.Errorf("write dirty file %q: %w", targetPath, err)
		}
	}

	return nil
}

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
