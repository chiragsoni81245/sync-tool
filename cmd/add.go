package cmd

import (
	"os"
	"slices"
    "log"
	"sync-tool/internal/db"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a sync target",
	Run: func(cmd *cobra.Command, args []string) {
        provider, _ := cmd.Flags().GetString("provider")
        mode, _ := cmd.Flags().GetString("mode")
        localDirectory, _ := cmd.Flags().GetString("local")
        remoteRef, _ := cmd.Flags().GetString("remote")

        if !slices.Contains([]db.SyncMode{db.ModePull, db.ModePush}, db.SyncMode(mode)) {
            log.Fatalf("invalid mode flag")
            return 
        }

        if !slices.Contains([]db.SyncProvider{db.ProviderGitHub, db.ProviderGDrive}, db.SyncProvider(provider)) {
            log.Fatalf("invalid provider flag")
            return 
        }

        // Validate local directory
        info, err := os.Stat(localDirectory) 

        if os.IsNotExist(err) {
            log.Fatalf("local path '%s' does not exists", localDirectory)
            return
        }

        if err != nil {
            log.Fatalf("unable to access '%s'", localDirectory)
            return
        } 

        if !info.IsDir() {
            log.Fatalf("'%s' is not a directory", localDirectory)
            return
        }


        // Save to DB
        syncTarget := db.SyncTarget{
            Mode:           db.SyncMode(mode),
            Provider:       db.SyncProvider(provider),
            LocalPath:      localDirectory,
            RemoteRef:      remoteRef,
            LastSyncStatus: db.StatusPending,
            StatusMessage:  "Waiting for first sync",
            LastSyncedAt:   nil,
        }

        if err := db.DB.Create(&syncTarget).Error; err != nil {
            log.Fatalf("Failed to save sync target: %v", err)
            return
        }

        log.Printf("Sync target added: %s", localDirectory)
    },
}


func init() {
    addCmd.Flags().String("provider", "", "Sync provider (github, gdrive)")
    addCmd.MarkFlagRequired("provider")

    addCmd.Flags().String("mode", "", "Sync mode (push, pull)")
    addCmd.MarkFlagRequired("mode")

    addCmd.Flags().String("local", "", "Local directory in/from which you want to sync")
    addCmd.MarkFlagRequired("local")

    addCmd.Flags().String("remote", "", "Remote directory reference in/from which you want to sync (e.g, github repo url, gdrive folder link)")
    addCmd.MarkFlagRequired("remote")
}
