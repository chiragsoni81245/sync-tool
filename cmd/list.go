// cmd/list.go
package cmd

import (
	"fmt"
	"text/tabwriter"
	"os"

	"github.com/spf13/cobra"
	"sync-tool/internal/db"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all sync targets",
	Run: func(cmd *cobra.Command, args []string) {
		var targets []db.SyncTarget
		if err := db.DB.Find(&targets).Error; err != nil {
			fmt.Println("Error fetching targets:", err)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tProvider\tMode\tPath\tRemoteRef\tLastStatus\tLastSynced")
		for _, t := range targets {
			lastSynced := "-"
			if t.LastSyncedAt != nil {
				lastSynced = t.LastSyncedAt.Format("2006-01-02 15:04:05")
			}
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t%s\n", t.ID, t.Provider, t.Mode, t.LocalPath, t.RemoteRef, t.LastSyncStatus, lastSynced)
		}
		w.Flush()
	},
}

