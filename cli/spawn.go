package cli

import (
	"fmt"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var (
	spawnDesc     string
	spawnParent   uint
	spawnPriority string
)

var spawnCmd = &cobra.Command{
	Use:   "spawn <title>",
	Short: "Create a new process",
	Long:  `Create a new process with an optional title, description, parent process, and priority.

Processes are modeled as OS processes with states: running, blocked, suspended, terminated.
Each process can have child sub-processes, creating an infinite hierarchy.`,
	Example: `  # Create a simple process
  taskctl spawn "Fix login bug"

  # Create with description and high priority
  taskctl spawn "Deploy to production" -d "Hotfix for critical issue" -P high

  # Create as child of another process
  taskctl spawn "Design database schema" -p 1`,

	Args: cobra.ExactArgs(1),
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
	spawnCmd.Flags().StringVarP(&spawnDesc, "desc", "D", "", "Process description (supports Markdown)")
	spawnCmd.Flags().UintVarP(&spawnParent, "parent", "p", 0, "Parent process ID (creates sub-process)")
	spawnCmd.Flags().StringVarP(&spawnPriority, "priority", "P", "medium", "Priority level: low, medium, high")
}
