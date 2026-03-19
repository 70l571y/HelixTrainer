// Package cfg содержит общие константы и конфигурацию приложения.
package cfg

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/jedib0t/go-pretty/v6/text"
)

const (
	// RepoURL - URL репозитория GitHub.
	RepoURL = "https://github.com/70l571y/HelixTrainer.git"

	// DBFileName - имя файла базы данных.
	DBFileName = "hxtrainer.db"
)

// GetConfigDir возвращает платформо-зависимую директорию конфигурации.
func GetConfigDir() string {
	appName := "hxtrainer"

	switch runtime.GOOS {
	case "windows":
		// Windows: %APPDATA%/hxtrainer
		appData := os.Getenv("APPDATA")
		if appData != "" {
			return filepath.Join(appData, appName)
		}
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "."+appName)
	default:
		// macOS / Linux: ~/.config/hxtrainer или XDG_CONFIG_HOME
		xdgConfig := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfig != "" {
			return filepath.Join(xdgConfig, appName)
		}
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".config", appName)
	}
}

// GetDBPath возвращает полный путь к файлу базы данных.
func GetDBPath() string {
	configDir := GetConfigDir()
	return filepath.Join(configDir, DBFileName)
}

// Стили для вывода
var (
	StyleSuccess = text.FgGreen
	StyleError   = text.FgRed
	StyleWarning = text.FgYellow
	StyleInfo    = text.FgCyan
	StyleBold    = text.Bold
)
