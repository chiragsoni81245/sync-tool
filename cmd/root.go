// cmd/root.go
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"sync-tool/internal/config"
	"sync-tool/internal/db"
	"sync-tool/internal/logger"
)

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
	config.LoadConfig("config.yaml")
	db.InitDB("data/sync.db")
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(listCmd)
}
