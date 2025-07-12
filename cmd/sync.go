// cmd/sync.go
package cmd

import (
	"fmt"

	"sync-tool/internal/db"
	"sync-tool/internal/provider"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Manually sync a target",
	Run: func(cmd *cobra.Command, args []string) {
        id, _ := cmd.Flags().GetInt("target")

		var target db.SyncTarget
		if err := db.DB.Where("id = ?", id).First(&target).Error; err != nil {
			fmt.Println("Error fetching target:", err)
			return
		}

        provider, _ := provider.GetProviderViaName(target.Provider)
        err := provider.Sync(target)
        if err != nil {
            fmt.Println("Error syncing target: ", err)
        }
	},
}


func init() {
    syncCmd.Flags().Int("target", -1, "Id of the target for sync")
    syncCmd.MarkFlagRequired("target")
}
