package cli

import (
	"fmt"
	"strconv"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var logUpdateCmd = &cobra.Command{
	Use:   "update <log_id> <new_content>",
	Short: "Update a log entry",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		logID, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			return err
		}

		newContent := args[1]

		if err := core.UpdateLog(uint(logID), newContent); err != nil {
			return err
		}

		fmt.Printf("Log %d updated\n", logID)
		return nil
	},
}

func init() {
	logCmd.AddCommand(logUpdateCmd)
}
