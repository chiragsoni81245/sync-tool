// internal/gitbhu/github.go
package github

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync-tool/internal/config"
	"sync-tool/internal/db"
	"sync-tool/internal/logger"
	"time"
)

type Github struct {
}

func New() *Github{
    return &Github{}
}

func (p *Github) PullSync(target db.SyncTarget) error {
    return nil
}

func (p *Github) PushSync(target db.SyncTarget) error {
	timeNow := time.Now()
	target.LastSyncedAt = &timeNow

	if err := initGitRepo(target.LocalPath); err != nil {
		target.LastSyncStatus = db.StatusFailed
		target.StatusMessage = "Git init error: " + err.Error()
		db.Save(&target)
		return nil
	}

	if err := configureGit(target.LocalPath); err != nil {
		target.LastSyncStatus = db.StatusFailed
		target.StatusMessage = "Git configure error: " + err.Error()
		db.Save(&target)
		return nil
	}

	if err := setRemote(target.LocalPath, target.RemoteRef); err != nil {
		target.LastSyncStatus = db.StatusFailed
		target.StatusMessage = "Set remote error: " + err.Error()
		db.Save(&target)
		return nil
	}

	if err := commitAndPush(target.LocalPath); err != nil {
		target.LastSyncStatus = db.StatusFailed
		target.StatusMessage = "Push error: " + err.Error()
		db.Save(&target)

        // If commit is not successful delete the remote
        if err := deleteRemote(target.LocalPath, target.RemoteRef); err != nil {
		    return nil
        }
		return nil
	}

	if err := deleteRemote(target.LocalPath, target.RemoteRef); err != nil {
		db.Save(&target)
		return nil
	}

	target.LastSyncStatus = db.StatusSuccess
	target.StatusMessage = "Sync successful"
	db.Save(&target)
	return nil
}


func buildAuthenticatedURL(repoURL string) string {
    return strings.Replace(repoURL, "https://github.com/", 
        fmt.Sprintf("https://%s:%s@github.com/", config.App.GitHubUsername, config.App.GitHubToken),
        1,
    )
}

func initGitRepo(path string) error {
	// Check if .git exists
	if _, err := os.Stat(fmt.Sprintf("%s/.git", path)); os.IsNotExist(err) {
		cmd := exec.Command("git", "init")
		cmd.Dir = path
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("git init failed: %v\n%s", err, out)
		}
		logger.Log.Infof("Initialized new git repo at %s", path)
	}
	return nil
}

func configureGit(path string) error {
    commands := [][]string{
        {"config", "user.email", config.App.GithubEmail},
        {"config", "user.name", config.App.GitHubUsername},
    }
    for _, args := range commands {
        cmd := exec.Command("git", args...)
        cmd.Dir = path
        if out, err := cmd.CombinedOutput(); err != nil {
            return fmt.Errorf("git %v failed: %v\n%s", args, err, out)
        }
    }

    logger.Log.Infof("Configured bot identity in git repo at %s", path)
	return nil
}

func setRemote(path, remoteURL string) error {
    remoteName := "github-auto-sync"
	cmd := exec.Command("git", "remote", "remove", remoteName)
	cmd.Dir = path
	cmd.Run() // ignore error if no remote exists

	cmd = exec.Command("git", "remote", "add", remoteName, buildAuthenticatedURL(remoteURL))
	cmd.Dir = path
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git remote add failed: %v\n%s", err, out)
	}
	logger.Log.Infof("Set remote %s: %s", remoteName, remoteURL)
	return nil
}

func deleteRemote(path, remoteURL string) error {
    remoteName := "github-auto-sync"
	cmd := exec.Command("git", "remote", "remove", remoteName)
	cmd.Dir = path
	cmd.Run() // ignore error if no remote exists

	logger.Log.Infof("Deleted remote %s: %s", remoteName, remoteURL)
	return nil
}

func commitAndPush(path string) error {
    remoteName := "github-auto-sync"
	commands := [][]string{
		{"add", "-A"},
		{"commit", "-m", "\"Auto sync\""},
		{"branch", "-M", "main"},
		{"push", "-u", remoteName, "main"},
	}
	for _, args := range commands {
		cmd := exec.Command("git", args...)
		cmd.Dir = path
		if out, err := cmd.CombinedOutput(); err != nil {
			// If commit fails because nothing changed, skip it
			if args[0] == "commit" && strings.Contains(string(out), "nothing to commit") {
				logger.Log.Infof("Nothing to commit in %s", path)
				continue
			}
			return fmt.Errorf("git %v failed: %v\n%s", args, err, out)
		}
	}
	logger.Log.Infof("Pushed changes for %s", path)
	return nil
}

