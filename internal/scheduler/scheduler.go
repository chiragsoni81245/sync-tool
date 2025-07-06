// internal/scheduler/scheduler.go
package scheduler

import (
	"github.com/robfig/cron/v3"
	"sync-tool/internal/config"
	"sync-tool/internal/db"
	"sync-tool/internal/provider"
	"sync-tool/internal/github"
	"sync-tool/internal/logger"
)

func Start() {
	c := cron.New()
	_, err := c.AddFunc(config.App.CronSchedule, func() {
		logger.Log.Infof("Running scheduled sync job")

		var targets []db.SyncTarget
		if err := db.DB.Find(&targets).Error; err != nil {
			logger.Log.Errorf("failed to fetch sync targets: %v", err)
			return
		}

		for _, target := range targets {
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
    }

    provider := providers[string(target.Provider)]
    switch target.Mode {
    case db.ModePull: provider.PullSync(target)
    case db.ModePush: provider.PushSync(target)
    }
}


