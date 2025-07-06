// cmd/root.go
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"sync-tool/internal/config"
	"sync-tool/internal/db"
	"sync-tool/internal/logger"
)

var configPath string
var dbPath string
var rootCmd = &cobra.Command{
	Use:   "sync-tool",
	Short: "A CLI tool to sync local directories with GitHub repos",
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
	config.LoadConfig(configPath)
	db.InitDB(dbPath)
}

func init() {
    rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to config file")
    rootCmd.MarkPersistentFlagRequired("config")
    rootCmd.PersistentFlags().StringVar(&dbPath, "db-path", "", "Path to database file")
    rootCmd.MarkPersistentFlagRequired("db-path")
    rootCmd.AddCommand(githubCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(listCmd)
}
