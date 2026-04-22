// Package cfg содержит общие константы и конфигурацию приложения.
package cfg

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/jedib0t/go-pretty/v6/text"
)

const (
	appName = "hxtrainer"

	// RepoURL - URL репозитория GitHub.
	RepoURL = "https://github.com/70l571y/HelixTrainer.git"

	// DBFileName - имя файла базы данных.
	DBFileName = "hxtrainer.db"

	// ChallengesDirName - имя директории с данными челленджей.
	ChallengesDirName = "challenges_data"
)

// GetConfigDir возвращает платформо-зависимую директорию конфигурации.
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

// GetDBPath возвращает полный путь к файлу базы данных.
func GetDBPath() string {
	configDir := GetConfigDir()
	return filepath.Join(configDir, DBFileName)
}

// GetChallengesRootDir возвращает корневую директорию встроенных челленджей в конфиге.
func GetChallengesRootDir() string {
	return filepath.Join(GetConfigDir(), ChallengesDirName)
}

// GetChallengesDir возвращает директорию с челленджами по умолчанию.
func GetChallengesDir() string {
	return filepath.Join(GetChallengesRootDir(), "go")
}

// Стили для вывода
var (
	StyleSuccess = text.FgGreen
	StyleError   = text.FgRed
	StyleWarning = text.FgYellow
	StyleInfo    = text.FgCyan
	StyleBold    = text.Bold
)
