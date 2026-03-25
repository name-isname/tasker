package cli

import (
	"fmt"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var (
	spawnDesc    string
	spawnParent  uint
	spawnPriority string
)

var spawnCmd = &cobra.Command{
	Use:   "spawn <title>",
	Short: "Create a new process",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		priority := core.PriorityMedium
		if spawnPriority != "" {
			priority = core.ProcessPriority(spawnPriority)
		}

		var parentID *uint
		if spawnParent > 0 {
			parentID = &spawnParent
		}

		process, err := core.CreateProcess(args[0], spawnDesc, parentID, priority)
		if err != nil {
			return err
		}
		fmt.Printf("Spawned process %d: %s\n", process.ID, process.Title)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(spawnCmd)
	spawnCmd.Flags().StringVarP(&spawnDesc, "desc", "d", "", "Process description")
	spawnCmd.Flags().UintVarP(&spawnParent, "parent", "p", 0, "Parent process ID")
	spawnCmd.Flags().StringVarP(&spawnPriority, "priority", "P", "medium", "Priority (low, medium, high)")
}
