// cmd/add.go
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"sync-tool/internal/db"
	"sync-tool/internal/logger"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a directory to be synced with a GitHub repo",
	Run: func(cmd *cobra.Command, args []string) {
		r := bufio.NewReader(os.Stdin)
		fmt.Print("Enter local directory path: ")
		dirPath, _ := r.ReadString('\n')
		dirPath = strings.TrimSpace(dirPath)

		fmt.Print("Enter GitHub repo URL: ")
		repoURL, _ := r.ReadString('\n')
		repoURL = strings.TrimSpace(repoURL)

		if dirPath == "" || repoURL == "" {
			logger.Log.Error("Directory path and repo URL are required")
			return
		}

		// Save to DB
		syncTarget := db.SyncTarget{
			Path:           dirPath,
			RepoURL:        repoURL,
			LastSyncStatus: db.StatusPending,
			StatusMessage:  "Waiting for first sync",
			LastSyncedAt:   nil,
		}

		if err := db.DB.Create(&syncTarget).Error; err != nil {
			logger.Log.Errorf("Failed to save sync target: %v", err)
			return
		}

		logger.Log.Infof("Sync target added: %s", dirPath)
	},
}

