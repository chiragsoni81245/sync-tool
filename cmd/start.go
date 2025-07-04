// cmd/start.go
package cmd

import (
	"github.com/spf13/cobra"
	"sync-tool/internal/scheduler"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start background cron job to sync directories",
	Run: func(cmd *cobra.Command, args []string) {
		scheduler.Start()
	},
}

