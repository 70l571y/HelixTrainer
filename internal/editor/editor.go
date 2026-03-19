// Package editor предоставляет функции для работы с редактором Helix.
package editor

import (
	"os"
	"os/exec"
	"path/filepath"
)

// OpenEditor открывает файл в редакторе Helix.
// Возвращает true если редактор успешно запущен.
func OpenEditor(filePath string, cwd string) bool {
	// Проверяем существование файла
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	// Получаем абсолютный путь
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return false
	}

	// Создаём команду для запуска hx
	cmd := exec.Command("hx", absPath)

	// Устанавливаем рабочую директорию
	if cwd != "" {
		cmd.Dir = cwd
	}

	// Подключаем стандартные потоки ввода/вывода
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Запускаем и ждём завершения
	err = cmd.Run()
	return err == nil
}

// HelixInstalled проверяет наличие редактора Helix в системе.
func HelixInstalled() bool {
	_, err := exec.LookPath("hx")
	return err == nil
}
