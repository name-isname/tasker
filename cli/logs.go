package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var (
	logsTail int
)

var logsCmd = &cobra.Command{
	Use:   "logs <pid>",
	Short: "Show logs for a process",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			return err
		}

		logEntries, err := core.GetLogs(uint(pid))
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(logEntries, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		if len(logEntries) == 0 {
			fmt.Printf("No logs for process %d\n", pid)
			return nil
		}

		// Apply tail limit
		if logsTail > 0 && len(logEntries) > logsTail {
			logEntries = logEntries[:logsTail]
		}

		// Get process info for header
		process, err := core.GetProcess(uint(pid))
		if err == nil {
			fmt.Printf("Logs for process #%d: %s\n", pid, process.Title)
			fmt.Println(strings.Repeat("-", 50))
		}

		// Display logs in reverse chronological order (already sorted)
		for _, log := range logEntries {
			icon := "📝"
			if log.LogType == core.LogTypeStateChange {
				icon = "🔄"
			}

			timestamp := log.CreatedAt.Format("15:04")
			fmt.Printf("%s [%s] %s\n", timestamp, icon, log.Content)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.Flags().IntVarP(&logsTail, "tail", "n", 0, "Show only N most recent logs")
}
