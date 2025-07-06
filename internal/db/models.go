// internal/db/models.go
package db

import (
	"sync-tool/internal/logger"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

func Save(t *SyncTarget) {
	if err := DB.Save(t).Error; err != nil {
		logger.Log.Errorf("Failed to update sync target: %v", err)
	}
}

type SyncMode string
const (
    ModePush SyncMode = "push"
    ModePull SyncMode = "pull"
)

type SyncProvider string
const (
    ProviderGitHub SyncProvider = "github"
    ProviderGDrive SyncProvider = "gdrive"
)

type SyncStatus string
const (
    StatusPending SyncStatus = "pending"
    StatusSuccess SyncStatus = "success"
    StatusFailed  SyncStatus = "failed"
)

type SyncTarget struct {
    ID             uint           `gorm:"primarykey"`
    Mode           SyncMode       `gorm:"index"`        // push or pull
    Provider       SyncProvider   `gorm:"index"`        // github or gdrive
    LocalPath      string
    RemoteRef      string         // repo URL or gdrive folder id
    LastSyncedAt   *time.Time
    LastSyncStatus SyncStatus
    StatusMessage  string
    CreatedAt      time.Time
    UpdatedAt      time.Time
}
