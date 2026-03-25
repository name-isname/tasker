package cli

import (
	"fmt"
	"strconv"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			return err
		}
		if err := core.DeleteTask(uint(id)); err != nil {
			return err
		}
		fmt.Printf("Task %d deleted\n", id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
