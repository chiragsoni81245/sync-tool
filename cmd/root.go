// cmd/root.go
package cmd

import (
	"os"

	"sync-tool/internal/config"
	"sync-tool/internal/db"
	"sync-tool/internal/logger"

	"github.com/spf13/cobra"
)

var configPath string
var rootCmd = &cobra.Command{
	Use:   "sync-tool",
	Short: "A CLI tool to sync local directories with GitHub, and Google drive",
}

func Execute() {
	cobra.OnInitialize(initApp)
	if err := rootCmd.Execute(); err != nil {
		logger.Log.Fatalf("Error: %v", err)
		os.Exit(1)
	}
}

func initApp() {
	logger.InitLogger()
    if configPath != "" {
	    config.LoadConfig(configPath)
	    db.InitDB(config.App.DatabaseFilepath)
    }
}

func init() {
    rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to config file")
    rootCmd.MarkPersistentFlagRequired("config")
    rootCmd.AddCommand(addCmd)
    rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(listCmd)
    rootCmd.AddCommand(syncCmd)
}
