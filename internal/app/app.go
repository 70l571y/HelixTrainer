// Package app содержит основные CLI команды приложения.
package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/70l571y/HelixTrainer/internal/cfg"
	"github.com/70l571y/HelixTrainer/internal/challenges"
	"github.com/70l571y/HelixTrainer/internal/database"
	"github.com/70l571y/HelixTrainer/internal/editor"
	"github.com/70l571y/HelixTrainer/internal/judge"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

var (
	challengesDir string
	verbose       bool
)

// InitCommands инициализирует все команды CLI.
func InitCommands(rootCmd *cobra.Command, cd string) {
	challengesDir = cd

	rootCmd.AddCommand(playCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(upgradeCmd)
}

// playCmd - команда для запуска челленджа.
var playCmd = &cobra.Command{
	Use:   "play [challenge_id]",
	Short: "Запустить челлендж",
	Long:  "Запускает сеанок прохождения челленджа. Если указан ID - запускает конкретный челлендж, иначе выбирает умно.",
	Run:   runPlay,
}

// listCmd - команда для списка челленджей.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Список всех челленджей",
	Long:  "Показывает все доступные челленджи с их статусом выполнения.",
	Run:   runList,
}

// statsCmd - команда для показа статистики.
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Показать статистику прогресса",
	Long:  "Показывает детальную статистику: последние попытки, лучшие времена, вехи.",
	Run:   runStats,
}

// upgradeCmd - команда для обновления.
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Обновить HelixTrainer",
	Long:  "Проверяет и устанавливает последнюю версию из GitHub.",
	Run:   runUpgrade,
}

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

	// Выбор челленджа
	var selected challenges.Challenge
	if len(args) > 0 && args[0] != "" {
		challengeID := args[0]
		for _, c := range challengeList {
			if c.ID == challengeID {
				selected = c
				break
			}
		}
		if selected.ID == "" {
			fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Челлендж с ID %s не найден!", challengeID))
			return
		}
	} else {
		selected = challenges.SelectSmartChallenge(challengeList)
	}

	fmt.Println(text.Bold.Sprintf("Запуск челленджа: %s", selected.ID))
	fmt.Println(selected.Description)

	// Подготовка временного файла
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

	// Чтение всех goal файлов
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

	// Инъекция заголовка/футера в start
	formattedStart := challenges.BuildFileContent(selected, string(startContent))

	for {
		// Создаём временную директорию
		tmpDir, err := os.MkdirTemp("", "helix_trainer*")
		if err != nil {
			fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Ошибка создания временной директории: %v", err))
			return
		}

		// Определяем имя основного файла
		ext := filepath.Ext(selected.StartPath)
		if ext == "" {
			ext = ".go"
		}
		mainFileName := "challenge" + ext
		tmpFilePath := filepath.Join(tmpDir, mainFileName)

		// Записываем файл
		if err := os.WriteFile(tmpFilePath, []byte(formattedStart), 0644); err != nil {
			fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Ошибка записи временного файла: %v", err))
			os.RemoveAll(tmpDir)
			return
		}

		// Копируем дополнительные файлы
		for _, extraPath := range selected.ExtraFilePaths {
			if _, err := os.Stat(extraPath); err == nil {
				destPath := filepath.Join(tmpDir, filepath.Base(extraPath))
				copyFile(extraPath, destPath)
			}
		}

		startTime := time.Now()

		// Открываем редактор
		if !editor.OpenEditor(tmpFilePath, tmpDir) {
			fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Ошибка запуска редактора Helix"))
			os.RemoveAll(tmpDir)
			return
		}

		endTime := time.Now()
		duration := endTime.Sub(startTime).Seconds()

		// Читаем отредактированный файл
		userContent, err := os.ReadFile(tmpFilePath)
		if err != nil {
			userContent = []byte{}
		}

		// Очистка временной директории
		os.RemoveAll(tmpDir)

		judgeMode := selected.JudgeMode
		if judgeMode == "" {
			judgeMode = "exact"
		}

		// Проверка решения
		isCorrect := false
		var matchingGoal string

		// Multi-file validation
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
			// Legacy single-file check
			for _, gContent := range goalContents {
				if judge.CheckSolution(string(userContent), gContent, judgeMode) {
					isCorrect = true
					matchingGoal = gContent
					break
				}
			}
		}

		// Проверка предыдущих попыток
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

		// Вехи
		ms := challenges.GetMilestone(duration, selected.AuthorTime)
		prevMsRank := 0
		if len(prevSuccesses) > 0 && selected.AuthorTime > 0 {
			prevMs := challenges.GetMilestone(prevBestTime, selected.AuthorTime)
			prevMsRank = prevMs.Rank
		}

		isNewMilestone := ms.Rank > prevMsRank

		// Логирование попытки
		database.LogAttempt(selected.ID, isCorrect, duration)

		// Вывод результата
		if isCorrect {
			displayFeedback(string(userContent), matchingGoal, judgeMode)

			msg := text.Bold.Sprintf(text.FgGreen.Sprintf("Вы выполнили задание за %.2fс", duration))
			if isNewRecord {
				msg += text.Bold.Sprintf(text.FgYellow.Sprintf(" 🏆 Новый рекорд!"))
			}
			if isNewMilestone {
				msg += text.Bold.Sprintf(text.FgMagenta.Sprintf("\n🎉 Новая веха: %s %s!", ms.Symbol, ms.Name))
			} else if ms.Name != "" {
				msg += fmt.Sprintf(" [%s %s]", ms.Symbol, ms.Name)
			}
			fmt.Println(msg)
		} else {
			displayFeedback(string(userContent), goalContents[0], judgeMode)
			fmt.Println(text.Bold.Sprintf(text.FgRed.Sprintf("❌ Не удалось. Время: %.2fс", duration)))
		}

		fmt.Println("\n[bold]Следующий челлендж: j | Повторить: k | Выход: esc/q[/bold]")

		// Ожидание ввода
		for {
			key := readKey()
			if key == "j" || key == "J" {
				fmt.Println(text.FgGreen.Sprintf("Загрузка следующего челленджа..."))
				runPlay(cmd, nil)
				return
			} else if key == "k" || key == "K" {
				// Повторить - выходим из внутреннего цикла
				break
			} else if key == "q" || key == "Q" || key == "\x1b" {
				fmt.Println("Выход.")
				return
			}
			// иначе игнорируем
		}
	}
}

func runList(cmd *cobra.Command, args []string) {
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
	challenges.SortForProgression(challengeList)

	// Создаём множество выполненных
	completedIDs := make(map[string]bool)
	for _, a := range attempts {
		if a.IsCorrect {
			completedIDs[a.ChallengeID] = true
		}
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Сложность", "Язык", "Метки", "Статус"})

	for _, c := range challengeList {
		status := "Не решено"
		displayID := c.ID

		if completedIDs[c.ID] {
			displayID = text.FgGreen.Sprintf(c.ID)
			status = text.Bold.Sprintf(text.FgGreen.Sprintf("Выполнено"))
		}

		labels := strings.Join(c.Tags, ", ")

		t.AppendRow(table.Row{displayID, c.Difficulty, c.Language, labels, status})
	}

	t.Render()
}

func runStats(cmd *cobra.Command, args []string) {
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

	if len(attempts) == 0 {
		fmt.Println(text.FgYellow.Sprintf("Попыток пока нет."))
		return
	}

	// Последние 20 попыток
	fmt.Println(text.Bold.Sprintf("Последняя активность (20)"))
	recentTable := table.NewWriter()
	recentTable.SetOutputMirror(os.Stdout)
	recentTable.AppendHeader(table.Row{"Время", "Челлендж", "Результат", "Время"})

	sort.Slice(attempts, func(i, j int) bool {
		return attempts[i].Timestamp.After(attempts[j].Timestamp)
	})

	for i := 0; i < 20 && i < len(attempts); i++ {
		a := attempts[i]
		result := text.FgGreen.Sprintf("Pass")
		if !a.IsCorrect {
			result = text.FgRed.Sprintf("Fail")
		}
		timeStr := a.Timestamp.Format("2006-01-02 15:04")
		recentTable.AppendRow(table.Row{timeStr, a.ChallengeID, result, fmt.Sprintf("%.2fс", a.Duration)})
	}
	recentTable.Render()
	fmt.Println()

	// Прогресс челленджей
	fmt.Println(text.Bold.Sprintf("Прогресс челленджей"))
	progressTable := table.NewWriter()
	progressTable.SetOutputMirror(os.Stdout)
	progressTable.AppendHeader(table.Row{"Челлендж", "Статус", "Лучшее время", "Веха", "Попыток"})

	// Группировка попыток по челленджам
	attemptsByID := make(map[string][]database.Attempt)
	for _, a := range attempts {
		attemptsByID[a.ChallengeID] = append(attemptsByID[a.ChallengeID], a)
	}

	challenges.SortForProgression(challengeList)

	for _, c := range challengeList {
		cAttempts := attemptsByID[c.ID]
		successfulAttempts := []database.Attempt{}
		for _, a := range cAttempts {
			if a.IsCorrect {
				successfulAttempts = append(successfulAttempts, a)
			}
		}

		isCompleted := len(successfulAttempts) > 0
		totalCount := len(cAttempts)

		status := "Не решено"
		if !isCompleted && totalCount > 0 {
			status = "Не решено"
		} else if !isCompleted {
			status = "Нет попыток"
		}

		bestTimeStr := "-"
		msDisplay := "-"

		if isCompleted {
			status = text.Bold.Sprintf(text.FgGreen.Sprintf("Выполнено"))
			bestTime := float64(1e9)
			for _, a := range successfulAttempts {
				if a.Duration < bestTime {
					bestTime = a.Duration
				}
			}
			bestTimeStr = text.FgGreen.Sprintf(fmt.Sprintf("%.2fс", bestTime))

			ms := challenges.GetMilestone(bestTime, c.AuthorTime)
			if ms.Name != "" {
				msDisplay = fmt.Sprintf("%s %s", ms.Symbol, ms.Name)
			}
		}

		progressTable.AppendRow(table.Row{c.ID, status, bestTimeStr, msDisplay, totalCount})
	}
	progressTable.Render()
}

func runUpgrade(cmd *cobra.Command, args []string) {
	fmt.Println(text.FgCyan.Sprintf("Проверка обновлений..."))

	latest := getLatestVersion()
	current := getCurrentVersion()

	if latest == "" {
		fmt.Fprintln(os.Stderr, text.FgRed.Sprintf("Не удалось получить информацию о последней версии."))
		return
	}

	if latest == current {
		fmt.Println(text.FgGreen.Sprintf("Уже установлена последняя версия (%s).", current))
		return
	}

	fmt.Println(text.FgYellow.Sprintf("Обновление с %s до %s...", current, latest))
	fmt.Println(text.FgYellow.Sprintf("Для обновления выполните:"))
	fmt.Println("  go install github.com/70l571y/HelixTrainer/cmd/hxtrainer@latest")
}

// Вспомогательные функции

func displayFeedback(userText, goalText, judgeMode string) {
	if judge.CheckSolution(userText, goalText, judgeMode) {
		fmt.Println("🎉 Успех! Решение верное.")
	} else {
		fmt.Println("❌ Решение неверное. Вот diff:")
		diff := judge.GenerateDiff(userText, goalText)
		fmt.Println(diff)
		if judgeMode == "ast" || judgeMode == "go_ast" {
			fmt.Println(text.FgYellow.Sprintf("Примечание: Режим AST. Структура должна совпадать точно, но форматирование может отличаться."))
		}
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

func readKey() string {
	// Простая реализация через bufio
	fmt.Print("")
	var input string
	fmt.Scanln(&input)
	if len(input) > 0 {
		return strings.ToLower(string(input[0]))
	}
	return ""
}

func getCurrentVersion() string {
	return "0.1.0" // TODO: получить из build flags
}

func getLatestVersion() string {
	repoBase := strings.TrimSuffix(cfg.RepoURL, ".git")
	rawURL := strings.Replace(repoBase, "github.com", "raw.githubusercontent.com", 1) + "/main/go/go.mod"

	resp, err := http.Get(rawURL)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	re := regexp.MustCompile(`module\s+\S+\n(?:.*\n)*?//\s*version:\s*(\S+)`)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) > 1 {
		return matches[1]
	}

	// Альтернативно: ищем в родительского проекта
	return ""
}
