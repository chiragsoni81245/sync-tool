// internal/git/git.go
package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync-tool/internal/logger"
    "sync-tool/internal/config"
)

func buildAuthenticatedURL(repoURL string) string {
    return strings.Replace(repoURL, "https://github.com/", 
        fmt.Sprintf("https://%s:%s@github.com/", config.App.GitHubUsername, config.App.GitHubToken),
        1,
    )
}

func InitGitRepo(path string) error {
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

func SetRemote(path, remoteURL string) error {
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

func DeleteRemote(path, remoteURL string) error {
    remoteName := "github-auto-sync"
	cmd := exec.Command("git", "remote", "remove", remoteName)
	cmd.Dir = path
	cmd.Run() // ignore error if no remote exists

	logger.Log.Infof("Deleted remote %s: %s", remoteName, remoteURL)
	return nil
}

func CommitAndPush(path string) error {
    remoteName := "github-auto-sync"
	commands := [][]string{
		{"add", "-A"},
		{"commit", "-m", "Auto sync"},
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

