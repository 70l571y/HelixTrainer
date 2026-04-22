// Package database предоставляет функции для работы с SQLite базой данных.
package database

import (
	"encoding/json"
	"io"
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
	dbMu sync.Mutex
)

// InitDB инициализирует базу данных и создаёт таблицы.
func InitDB() error {
	dbMu.Lock()
	defer dbMu.Unlock()

	if db != nil {
		return nil
	}

	dbPath := cfg.GetDBPath()

	if err := ensureDir(dbPath); err != nil {
		return err
	}

	openedDB, err := gorm.Open(sqlite.Open("file:"+dbPath), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := openedDB.AutoMigrate(&Attempt{}); err != nil {
		return err
	}

	db = openedDB
	return nil
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

// ResetAttempts удаляет все записи о попытках.
func ResetAttempts() (int64, error) {
	if err := InitDB(); err != nil {
		return 0, err
	}

	result := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Attempt{})
	return result.RowsAffected, result.Error
}

func ExportAttempts(w io.Writer) error {
	attempts, err := GetAllAttempts()
	if err != nil {
		return err
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(attempts)
}

func ImportAttempts(r io.Reader, replace bool) (int, error) {
	if err := InitDB(); err != nil {
		return 0, err
	}

	var attempts []Attempt
	if err := json.NewDecoder(r).Decode(&attempts); err != nil {
		return 0, err
	}

	if replace {
		if _, err := ResetAttempts(); err != nil {
			return 0, err
		}
	}

	for i := range attempts {
		attempts[i].ID = 0
	}

	if len(attempts) == 0 {
		return 0, nil
	}

	if err := db.Create(&attempts).Error; err != nil {
		return 0, err
	}

	return len(attempts), nil
}
