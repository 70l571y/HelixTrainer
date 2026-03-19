// Package database предоставляет функции для работы с SQLite базой данных.
package database

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/70l571y/HelixTrainer/internal/cfg"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Attempt представляет запись о попытке прохождения челленджа.
type Attempt struct {
	ID          uint      `gorm:"primaryKey"`
	ChallengeID string    `gorm:"index"`
	Timestamp   time.Time `gorm:"autoCreateTime"`
	IsCorrect   bool
	Duration    float64
}

var (
	db   *gorm.DB
	once sync.Once
)

// InitDB инициализирует базу данных и создаёт таблицы.
func InitDB() error {
	var initErr error
	once.Do(func() {
		dbPath := cfg.GetDBPath()

		// Создаём директорию если не существует
		var err error
		if err = ensureDir(dbPath); err != nil {
			initErr = err
			return
		}

		dsn := "file:" + dbPath
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
		if err != nil {
			initErr = err
			return
		}

		initErr = db.AutoMigrate(&Attempt{})
	})
	return initErr
}

// ensureDir создаёт директорию для файла если она не существует.
func ensureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	if dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0755)
}

// LogAttempt записывает попытку прохождения челленджа в базу данных.
func LogAttempt(challengeID string, isCorrect bool, duration float64) (*Attempt, error) {
	if err := InitDB(); err != nil {
		return nil, err
	}

	attempt := Attempt{
		ChallengeID: challengeID,
		IsCorrect:   isCorrect,
		Duration:    duration,
	}

	result := db.Create(&attempt)
	return &attempt, result.Error
}

// GetAttempts возвращает все попытки или попытки для конкретного челленджа.
func GetAttempts(challengeID ...string) ([]Attempt, error) {
	if err := InitDB(); err != nil {
		return nil, err
	}

	var attempts []Attempt
	query := db.Order("timestamp DESC")

	if len(challengeID) > 0 && challengeID[0] != "" {
		query = query.Where("challenge_id = ?", challengeID[0])
	}

	result := query.Find(&attempts)
	return attempts, result.Error
}

// GetAttemptsByChallenge возвращает попытки для конкретного челленджа.
func GetAttemptsByChallenge(challengeID string) ([]Attempt, error) {
	return GetAttempts(challengeID)
}

// GetAllAttempts возвращает все попытки.
func GetAllAttempts() ([]Attempt, error) {
	return GetAttempts("")
}
