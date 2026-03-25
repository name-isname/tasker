package cli

import (
	"fmt"
	"strconv"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status <id> <new-status>",
	Short: "Change the status of a process",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			return err
		}

		newStatus := core.ProcessStatus(args[1])
		if !isValidStatus(newStatus) {
			return fmt.Errorf("invalid status: %s (must be: running, blocked, suspended, terminated)", args[1])
		}

		if err := core.SetProcessStatus(uint(id), newStatus); err != nil {
			return err
		}
		fmt.Printf("Process %d status changed to %s\n", id, newStatus)
		return nil
	},
}

func isValidStatus(status core.ProcessStatus) bool {
	switch status {
	case core.StatusRunning, core.StatusBlocked, core.StatusSuspended, core.StatusTerminated:
		return true
	default:
		return false
	}
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
