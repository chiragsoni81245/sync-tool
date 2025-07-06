// cmd/delete.go
package cmd

import (
	"log"
	"sync-tool/internal/db"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a sync target",
	Run: func(cmd *cobra.Command, args []string) {
        id, _ := cmd.Flags().GetInt("id")

        if id <= 0 {
            log.Fatal("Invalid target id")
            return
        }

		var target db.SyncTarget
        db.DB.Where("id = ?", id).First(&target)
        db.DB.Delete(target)
	},
}

func init() {
    deleteCmd.Flags().Int("id", -1, "Id of the sync target you want to delete")
    deleteCmd.MarkFlagRequired("id")
}
