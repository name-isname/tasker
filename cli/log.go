package cli

import (
	"fmt"
	"strconv"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log <pid> <content>",
	Short: "Append a log entry to a process",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			return err
		}

		content := args[1]

		log, err := core.AddLog(uint(pid), core.LogTypeProgress, content)
		if err != nil {
			return err
		}

		fmt.Printf("Log #%d added to process %d\n", log.ID, pid)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
