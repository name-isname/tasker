package cli

import (
	"fmt"
	"strconv"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill <pid>",
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
		fmt.Printf("Killed process %d\n", id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(killCmd)
}
