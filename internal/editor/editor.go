// Package editor предоставляет функции для работы с редактором Helix.
package editor

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var ErrHelixNotFound = errors.New("helix executable not found")

// OpenEditor открывает файл в редакторе Helix.
func OpenEditor(filePath string, cwd string) error {
	// Проверяем существование файла
	if _, err := os.Stat(filePath); err != nil {
		return err
	}

	// Получаем абсолютный путь
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	// Создаём команду для запуска hx
	path, err := exec.LookPath("hx")
	if err != nil {
		return fmt.Errorf("%w: install Helix and make sure hx is in PATH", ErrHelixNotFound)
	}

	cmd := exec.Command(path, absPath)

	// Устанавливаем рабочую директорию
	if cwd != "" {
		cmd.Dir = cwd
	}

	// Подключаем стандартные потоки ввода/вывода
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Запускаем и ждём завершения
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run hx: %w", err)
	}

	return nil
}

// HelixInstalled проверяет наличие редактора Helix в системе.
func HelixInstalled() bool {
	_, err := exec.LookPath("hx")
	return err == nil
}
