package cli

import (
	"fmt"
	"strconv"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:   "complete <id>",
	Short: "Mark a task as completed",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			return err
		}
		if err := core.CompleteTask(uint(id)); err != nil {
			return err
		}
		fmt.Printf("Task %d marked as completed\n", id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)
}
