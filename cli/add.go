package cli

import (
	"fmt"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <title>",
	Short: "Add a new task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		task, err := core.AddTask(args[0])
		if err != nil {
			return err
		}
		fmt.Printf("Task added: %s (ID: %d)\n", task.Title, task.ID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
