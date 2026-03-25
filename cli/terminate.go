package cli

import (
	"fmt"
	"strconv"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var (
	terminateReason string
)

var terminateCmd = &cobra.Command{
	Use:   "terminate <pid>",
	Short: "Terminate a process (mark as completed)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			return err
		}

		if err := core.ChangeProcessState(uint(id), core.StatusTerminated, terminateReason); err != nil {
			return err
		}

		if terminateReason != "" {
			fmt.Printf("Process %d terminated: %s\n", id, terminateReason)
		} else {
			fmt.Printf("Process %d terminated\n", id)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(terminateCmd)
	terminateCmd.Flags().StringVarP(&terminateReason, "message", "m", "", "Completion message")
}
