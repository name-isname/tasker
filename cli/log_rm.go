package cli

import (
	"fmt"
	"strconv"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var logRmCmd = &cobra.Command{
	Use:   "rm <log_id>",
	Short: "Delete a log entry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		logID, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			return err
		}

		if err := core.DeleteLog(uint(logID)); err != nil {
			return err
		}

		fmt.Printf("Log %d deleted\n", logID)
		return nil
	},
}

func init() {
	logCmd.AddCommand(logRmCmd)
}
