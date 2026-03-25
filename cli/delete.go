package cli

import (
	"fmt"
	"strconv"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a process (and all its children)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			return err
		}
		if err := core.DeleteProcess(uint(id)); err != nil {
			return err
		}
		fmt.Printf("Process %d deleted\n", id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
