package cli

import (
	"fmt"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var (
	addDescription string
	addParentID    uint
	addPriority    string
)

var addCmd = &cobra.Command{
	Use:   "add <title>",
	Short: "Add a new process",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		priority := core.PriorityMedium
		if addPriority != "" {
			priority = core.ProcessPriority(addPriority)
		}

		var parentID *uint
		if addParentID > 0 {
			parentID = &addParentID
		}

		process, err := core.CreateProcess(args[0], addDescription, parentID, priority)
		if err != nil {
			return err
		}
		fmt.Printf("Process created: %s (ID: %d, Status: %s)\n", process.Title, process.ID, process.Status)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&addDescription, "description", "d", "", "Process description")
	addCmd.Flags().UintVarP(&addParentID, "parent", "p", 0, "Parent process ID")
	addCmd.Flags().StringVarP(&addPriority, "priority", "P", "medium", "Priority (low, medium, high)")
}
