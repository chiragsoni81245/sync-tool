// internal/db/models.go
package db

import (
	"sync-tool/internal/logger"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	StatusSuccess = "success"
	StatusFailed  = "failed"
	StatusPending = "pending"
)

var DB *gorm.DB

func InitDB(path string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		logger.Log.Fatalf("failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(&SyncTarget{})
	if err != nil {
		logger.Log.Fatalf("failed to run migrations: %v", err)
	}
}

type SyncTarget struct {
	ID              uint           `gorm:"primaryKey"`
	Path            string         `gorm:"unique;not null"` // Local directory path
	RepoURL         string         `gorm:"not null"`         // GitHub repository URL
	LastSyncedAt    *time.Time
	LastSyncStatus  string         `gorm:"not null"`         // One of: "success", "failed", "pending"
	StatusMessage   string         `gorm:"type:text"`        // e.g. "Sync successful" or error message
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

