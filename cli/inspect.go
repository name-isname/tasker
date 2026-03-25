package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect <pid>",
	Short: "Show detailed information about a process",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			return err
		}

		process, err := core.GetProcess(uint(id))
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(process, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		// Pretty print process details
		fmt.Printf("┌─────────────────────────────────────────────────┐\n")
		fmt.Printf("│ Process #%d%-37s │\n", process.ID, "")
		fmt.Printf("├─────────────────────────────────────────────────┤\n")
		fmt.Printf("│ Title:       %-37s │\n", truncate(process.Title, 37))
		fmt.Printf("│ Status:      %-37s │\n", process.Status)
		fmt.Printf("│ Priority:    %-37s │\n", process.Priority)
		fmt.Printf("│ Ranking:     %-37.2f │\n", process.Ranking)

		if process.ParentID != nil {
			fmt.Printf("│ Parent:      #%d%-35s │\n", *process.ParentID, "")
		} else {
			fmt.Printf("│ Parent:      %-37s │\n", "(none)")
		}

		fmt.Printf("│ Created:     %-37s │\n", process.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("│ Updated:     %-37s │\n", process.UpdatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("├─────────────────────────────────────────────────┤\n")

		if process.Description != "" {
			// Word wrap description
			lines := wordWrap(process.Description, 51)
			for i, line := range lines {
				if i == 0 {
					fmt.Printf("│ Description: %-36s │\n", line)
				} else {
					fmt.Printf("│             %-36s │\n", line)
				}
			}
			fmt.Printf("├─────────────────────────────────────────────────┤\n")
		}

		// Show log count
		logCount := len(process.Logs)
		fmt.Printf("│ Logs:        %d entries%-28s │\n", logCount, "")
		fmt.Printf("└─────────────────────────────────────────────────┘\n")

		return nil
	},
}

func truncate(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s + strings.Repeat(" ", maxLen-len(s))
}

func wordWrap(s string, width int) []string {
	if len(s) <= width {
		return []string{s}
	}

	var lines []string
	for len(s) > width {
		lines = append(lines, s[:width])
		s = s[width:]
	}
	if len(s) > 0 {
		lines = append(lines, s)
	}
	return lines
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	inspectCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
}
