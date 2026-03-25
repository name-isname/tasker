package cli

import (
	"fmt"
	"strconv"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var (
	statusReason string
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

		// Use ChangeProcessState to atomically update status and create log
		if err := core.ChangeProcessState(uint(id), newStatus, statusReason); err != nil {
			return err
		}

		if statusReason != "" {
			fmt.Printf("Process %d status changed to %s (Reason: %s)\n", id, newStatus, statusReason)
		} else {
			fmt.Printf("Process %d status changed to %s\n", id, newStatus)
		}
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
	statusCmd.Flags().StringVarP(&statusReason, "reason", "r", "", "Reason for status change")
}
