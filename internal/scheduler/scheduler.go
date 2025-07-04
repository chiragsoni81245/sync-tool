// internal/scheduler/scheduler.go
package scheduler

import (
	"time"

	"github.com/robfig/cron/v3"
	"sync-tool/internal/config"
	"sync-tool/internal/db"
	"sync-tool/internal/git"
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
	timeNow := time.Now()
	target.LastSyncedAt = &timeNow

	if err := git.InitGitRepo(target.Path); err != nil {
		target.LastSyncStatus = db.StatusFailed
		target.StatusMessage = "Git init error: " + err.Error()
		save(&target)
		return
	}

	if err := git.ConfigureGit(target.Path); err != nil {
		target.LastSyncStatus = db.StatusFailed
		target.StatusMessage = "Git configure error: " + err.Error()
		save(&target)
		return
	}

	if err := git.SetRemote(target.Path, target.RepoURL); err != nil {
		target.LastSyncStatus = db.StatusFailed
		target.StatusMessage = "Set remote error: " + err.Error()
		save(&target)
		return
	}

	if err := git.CommitAndPush(target.Path); err != nil {
		target.LastSyncStatus = db.StatusFailed
		target.StatusMessage = "Push error: " + err.Error()
		save(&target)

        // If commit is not successful delete the remote
        if err := git.DeleteRemote(target.Path, target.RepoURL); err != nil {
            return
        }
		return
	}

	if err := git.DeleteRemote(target.Path, target.RepoURL); err != nil {
		save(&target)
		return
	}

	target.LastSyncStatus = db.StatusSuccess
	target.StatusMessage = "Sync successful"
	save(&target)
}

func save(t *db.SyncTarget) {
	if err := db.DB.Save(t).Error; err != nil {
		logger.Log.Errorf("Failed to update sync target: %v", err)
	}
}

