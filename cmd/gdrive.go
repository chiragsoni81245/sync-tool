// cmd/gdrive.go
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


var gdriveCmd = &cobra.Command{
	Use:   "gdrive",
	Short: "This gdrive command will allow user to do pull sync of a gdrive folder",
}

var gdrivePullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Add a directory to be synced with a Gdrive folder via pulling new changes periodically",
	Run: func(cmd *cobra.Command, args []string) {
		r := bufio.NewReader(os.Stdin)
		fmt.Print("Enter local directory path: ")
		dirPath, _ := r.ReadString('\n')
		dirPath = strings.TrimSpace(dirPath)

		fmt.Print("Enter Google Drive folder URL: ")
		remoteURL, _ := r.ReadString('\n')
		remoteURL = strings.TrimSpace(remoteURL)

		if dirPath == "" || remoteURL == "" {
			logger.Log.Error("Directory path and Google Drive folder URL are required")
			return
		}

		// Save to DB
        syncTarget := db.SyncTarget{
            Mode:           db.ModePull,
            Provider:       db.ProviderGDrive,
            LocalPath:      dirPath,
            RemoteRef:      remoteURL,
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


func init() {
	gdriveCmd.AddCommand(gdrivePullCmd)
}
