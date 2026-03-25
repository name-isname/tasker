package cli

import (
	"fmt"
	"strconv"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var (
	wakeReason string
)

var wakeCmd = &cobra.Command{
	Use:   "wake <pid>",
	Short: "Wake up a process (set status back to running)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			return err
		}

		if err := core.ChangeProcessState(uint(id), core.StatusRunning, wakeReason); err != nil {
			return err
		}

		if wakeReason != "" {
			fmt.Printf("Process %d woke up: %s\n", id, wakeReason)
		} else {
			fmt.Printf("Process %d woke up\n", id)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(wakeCmd)
	wakeCmd.Flags().StringVarP(&wakeReason, "message", "m", "", "Reason for waking up")
}
