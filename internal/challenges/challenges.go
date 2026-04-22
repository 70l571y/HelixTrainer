// Package challenges предоставляет функции для загрузки и управления челленджами.
package challenges

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/70l571y/HelixTrainer/internal/database"
)

// Challenge представляет структуру данных челленджа.
type Challenge struct {
	ID             string            `json:"id"`
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	Difficulty     string            `json:"difficulty"`
	Language       string            `json:"language"`
	JudgeMode      string            `json:"judge_mode"`
	StartFile      string            `json:"start_file"`
	MainFileName   string            `json:"main_file_name"`
	GoalFile       string            `json:"goal_file"`
	GoalFiles      []string          `json:"goal_files"`
	Tips           string            `json:"tips"`
	Tags           []string          `json:"tags"`
	AuthorTime     float64           `json:"author_time"`
	ExtraFiles     []string          `json:"extra_files"`
	Validation     map[string]string `json:"validation"`
	GitDirtyFiles  map[string]string `json:"git_dirty_files"`
	DirPath        string            `json:"-"`
	StartPath      string            `json:"-"`
	GoalPaths      []string          `json:"-"`
	ExtraFilePaths []string          `json:"-"`
	ValidationMap  map[string]string `json:"-"`
}

type ChallengeFilter struct {
	Difficulty string
	Tags       []string
	Track      string
}

// GetCommentPrefix возвращает префикс комментариев для языка.
func GetCommentPrefix(language string) string {
	switch language {
	case "rust", "c", "cpp", "java", "javascript", "typescript", "go":
		return "//"
	case "sql":
		return "--"
	default:
		return "#"
	}
}

// BuildFileContent добавляет заголовок и подсказки к содержимому файла.
func BuildFileContent(challenge Challenge, content string) string {
	commentPrefix := GetCommentPrefix(challenge.Language)

	var sb strings.Builder

	// Заголовок
	sb.WriteString(commentPrefix + " " + challenge.ID + "\n")
	sb.WriteString(commentPrefix + " Task: " + challenge.Description + "\n\n")

	// Содержимое
	sb.WriteString(content)

	// Подсказки (footer)
	if challenge.Tips != "" {
		sb.WriteString("\n\n")
		lines := strings.Split(challenge.Tips, "\n")
		for _, line := range lines {
			sb.WriteString(commentPrefix + " " + line + "\n")
		}
	}

	return sb.String()
}

// LoadChallenges загружает все челленджи из директории challenges_data.
func LoadChallenges(challengesDir string) ([]Challenge, error) {
	var challenges []Challenge

	err := filepath.Walk(challengesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == "config.json" {
			challenge, err := loadChallengeFromFile(path)
			if err != nil {
				return err
			}
			challenges = append(challenges, challenge)
		}

		return nil
	})

	return challenges, err
}

// loadChallengeFromFile загружает челлендж из файла config.json.
func loadChallengeFromFile(configPath string) (Challenge, error) {
	var challenge Challenge

	data, err := os.ReadFile(configPath)
	if err != nil {
		return challenge, err
	}

	if err := json.Unmarshal(data, &challenge); err != nil {
		return challenge, err
	}

	// Разрешаем пути относительно файла конфигурации
	dirPath := filepath.Dir(configPath)
	challenge.DirPath = dirPath
	challenge.StartPath = filepath.Join(dirPath, challenge.StartFile)

	// Обработка goal файлов
	if len(challenge.GoalFiles) == 0 && challenge.GoalFile != "" {
		challenge.GoalFiles = []string{challenge.GoalFile}
	}

	for _, gf := range challenge.GoalFiles {
		challenge.GoalPaths = append(challenge.GoalPaths, filepath.Join(dirPath, gf))
	}

	// Обработка дополнительных файлов
	for _, ef := range challenge.ExtraFiles {
		challenge.ExtraFilePaths = append(challenge.ExtraFilePaths, filepath.Join(dirPath, ef))
	}

	// Обработка validation map
	if len(challenge.Validation) > 0 {
		challenge.ValidationMap = make(map[string]string)
		for filename, goalFilename := range challenge.Validation {
			challenge.ValidationMap[filename] = filepath.Join(dirPath, goalFilename)
		}
	}

	return challenge, nil
}

// Milestone представляет информацию о вехе.
type Milestone struct {
	Name   string
	Symbol string
	Rank   int
}

// GetMilestone возвращает веху на основе времени и авторского времени.
func GetMilestone(duration, authorTime float64) Milestone {
	if authorTime <= 0 {
		return Milestone{}
	}

	// Округляем до 2 знаков
	duration = float64(int(duration*100)) / 100

	if duration <= authorTime {
		return Milestone{Name: "Author", Symbol: "🟢", Rank: 4}
	} else if duration <= authorTime*1.25 {
		return Milestone{Name: "Gold", Symbol: "🥇", Rank: 3}
	} else if duration <= authorTime*1.75 {
		return Milestone{Name: "Silver", Symbol: "🥈", Rank: 2}
	} else if duration <= authorTime*3.0 {
		return Milestone{Name: "Bronze", Symbol: "🥉", Rank: 1}
	}

	return Milestone{}
}

func difficultyRank(difficulty string) int {
	switch strings.ToLower(strings.TrimSpace(difficulty)) {
	case "easy":
		return 0
	case "medium":
		return 1
	case "hard":
		return 2
	default:
		return 3
	}
}

func lessByProgression(a, b Challenge) bool {
	aRank := difficultyRank(a.Difficulty)
	bRank := difficultyRank(b.Difficulty)
	if aRank != bRank {
		return aRank < bRank
	}

	if a.ID != b.ID {
		return a.ID < b.ID
	}

	return a.Title < b.Title
}

// SortForProgression сортирует challenge-ы в порядке обучения:
// сначала Easy, затем Medium, затем Hard, внутри уровня по ID.
func SortForProgression(challenges []Challenge) {
	sort.Slice(challenges, func(i, j int) bool {
		return lessByProgression(challenges[i], challenges[j])
	})
}

func FilterChallenges(challenges []Challenge, filter ChallengeFilter) []Challenge {
	var filtered []Challenge
	wantDifficulty := strings.ToLower(strings.TrimSpace(filter.Difficulty))
	wantTrack := normalizeTrackTag(filter.Track)
	wantTags := make([]string, 0, len(filter.Tags))
	for _, tag := range filter.Tags {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			wantTags = append(wantTags, tag)
		}
	}

	for _, challenge := range challenges {
		if wantDifficulty != "" && strings.ToLower(strings.TrimSpace(challenge.Difficulty)) != wantDifficulty {
			continue
		}
		if wantTrack != "" && !hasTag(challenge.Tags, wantTrack) {
			continue
		}
		if len(wantTags) > 0 && !hasAnyTag(challenge.Tags, wantTags) {
			continue
		}
		filtered = append(filtered, challenge)
	}

	return filtered
}

func normalizeTrackTag(track string) string {
	switch strings.ToLower(strings.TrimSpace(track)) {
	case "", "any":
		return ""
	case "core", "hotkey", "core_hotkey":
		return "track_core_hotkey"
	case "optional", "command", "command_line", "optional_command_line":
		return "track_optional_command_line"
	default:
		return track
	}
}

func hasTag(tags []string, want string) bool {
	for _, tag := range tags {
		if tag == want {
			return true
		}
	}
	return false
}

func hasAnyTag(tags []string, want []string) bool {
	for _, tag := range want {
		if hasTag(tags, tag) {
			return true
		}
	}
	return false
}

// SelectSmartChallenge выбирает следующий челлендж на основе истории.
func SelectSmartChallenge(challenges []Challenge) Challenge {
	if len(challenges) == 0 {
		return Challenge{}
	}

	// Получаем все попытки
	allAttempts, _ := database.GetAllAttempts()

	// Создаём карту попыток по челленджам
	attemptsMap := make(map[string][]database.Attempt)
	for _, c := range challenges {
		attemptsMap[c.ID] = []database.Attempt{}
	}
	for _, attempt := range allAttempts {
		if _, ok := attemptsMap[attempt.ChallengeID]; ok {
			attemptsMap[attempt.ChallengeID] = append(attemptsMap[attempt.ChallengeID], attempt)
		}
	}

	var neverAttempted []Challenge
	var failedLast []Challenge
	var solved []Challenge

	for _, challenge := range challenges {
		cAttempts := attemptsMap[challenge.ID]
		if len(cAttempts) == 0 {
			neverAttempted = append(neverAttempted, challenge)
			continue
		}

		// Сортируем по времени (новые первые)
		sort.Slice(cAttempts, func(i, j int) bool {
			return cAttempts[i].Timestamp.After(cAttempts[j].Timestamp)
		})

		lastAttempt := cAttempts[0]
		if !lastAttempt.IsCorrect {
			failedLast = append(failedLast, challenge)
		} else {
			solved = append(solved, challenge)
		}
	}

	// Приоритет: никогда не пробованные
	if len(neverAttempted) > 0 {
		SortForProgression(neverAttempted)
		return neverAttempted[0]
	}

	// Затем: последние неудачные
	if len(failedLast) > 0 {
		SortForProgression(failedLast)
		return failedLast[0]
	}

	// Затем: решённые, тоже в порядке progression
	if len(solved) > 0 {
		SortForProgression(solved)
		return solved[0]
	}

	// fallback
	SortForProgression(challenges)
	return challenges[0]
}

func SelectProgressionChallenge(challenges []Challenge, attempts []database.Attempt) Challenge {
	if len(challenges) == 0 {
		return Challenge{}
	}

	filtered := append([]Challenge(nil), challenges...)
	SortForProgression(filtered)
	solved := solvedChallenges(attempts)
	for _, challenge := range filtered {
		if !solved[challenge.ID] {
			return challenge
		}
	}
	return filtered[0]
}

func SelectWeakestChallenge(challenges []Challenge, attempts []database.Attempt) Challenge {
	if len(challenges) == 0 {
		return Challenge{}
	}

	solved := solvedChallenges(attempts)
	tagScores := make(map[string]int)
	tagPending := make(map[string]int)
	tagSolved := make(map[string]int)

	for _, challenge := range challenges {
		for _, tag := range nonTrackTags(challenge.Tags) {
			if solved[challenge.ID] {
				tagSolved[tag]++
			} else {
				tagPending[tag]++
			}
		}
	}

	for _, attempt := range attempts {
		if attempt.IsCorrect {
			continue
		}
		challenge := findChallengeByID(challenges, attempt.ChallengeID)
		if challenge.ID == "" {
			continue
		}
		for _, tag := range nonTrackTags(challenge.Tags) {
			tagScores[tag] += 2
		}
	}

	best := Challenge{}
	bestScore := -1
	for _, challenge := range challenges {
		if solved[challenge.ID] {
			continue
		}

		score := 0
		for _, tag := range nonTrackTags(challenge.Tags) {
			score += tagScores[tag] + tagPending[tag]
			if tagSolved[tag] == 0 {
				score += 5
			}
		}
		if score > bestScore || (score == bestScore && lessByProgression(challenge, best)) {
			best = challenge
			bestScore = score
		}
	}

	if best.ID != "" {
		return best
	}
	return SelectProgressionChallenge(challenges, attempts)
}

func solvedChallenges(attempts []database.Attempt) map[string]bool {
	solved := make(map[string]bool)
	for _, attempt := range attempts {
		if attempt.IsCorrect {
			solved[attempt.ChallengeID] = true
		}
	}
	return solved
}

func nonTrackTags(tags []string) []string {
	var out []string
	for _, tag := range tags {
		if strings.HasPrefix(tag, "track_") {
			continue
		}
		out = append(out, tag)
	}
	return out
}

func findChallengeByID(challenges []Challenge, id string) Challenge {
	for _, challenge := range challenges {
		if challenge.ID == id {
			return challenge
		}
	}
	return Challenge{}
}
