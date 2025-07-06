// internal/scheduler/scheduler.go
package scheduler

import (
	"sync-tool/internal/config"
	"sync-tool/internal/db"
	"sync-tool/internal/gdrive"
	"sync-tool/internal/github"
	"sync-tool/internal/logger"
	"sync-tool/internal/provider"

	"github.com/robfig/cron/v3"
)

func Start() {
	c := cron.New()
	_, err := c.AddFunc(config.App.CronSchedule, func() {
		logger.Log.Infof("Running scheduled sync job")

		var pullSyncTargets []db.SyncTarget
		if err := db.DB.Where("mode = ?", db.ModePull).Find(&pullSyncTargets).Error; err != nil {
			logger.Log.Errorf("failed to fetch pull sync targets: %v", err)
			return
		}
		for _, target := range pullSyncTargets {
			syncOne(target)
		}

		var pushSyncTargets []db.SyncTarget
		if err := db.DB.Where("mode = ?", db.ModePush).Find(&pushSyncTargets).Error; err != nil {
			logger.Log.Errorf("failed to fetch push sync targets: %v", err)
			return
		}
		for _, target := range pushSyncTargets {
			syncOne(target)
		}
	})
	if err != nil {
		logger.Log.Fatalf("failed to schedule sync job: %v", err)
	}

	c.Start()
    logger.Log.Infof("Started cron scheduler")
	select {} // Block forever
}

func syncOne(target db.SyncTarget) {
    providers := map[string]provider.Provider{
        string(db.ProviderGitHub): github.New(),
        string(db.ProviderGDrive): gdrive.New(),
    }

    provider := providers[string(target.Provider)]
    switch target.Mode {
    case db.ModePull: provider.PullSync(target)
    case db.ModePush: provider.PushSync(target)
    }
}


