package cli

import (
	"encoding/json"
	"fmt"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
	psStatus   string
	psTree     bool
)

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List all processes",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Tree mode
		if psTree {
			treeOutput, err := core.FormatFullTree()
			if err != nil {
				return err
			}
			fmt.Print(treeOutput)
			return nil
		}

		// List mode
		var status *core.ProcessStatus
		if psStatus != "" && psStatus != "all" {
			s := core.ProcessStatus(psStatus)
			status = &s
		}

		processes, err := core.ListProcesses(status)
		if err != nil {
			return err
		}

		if jsonOutput {
			data, err := json.Marshal(processes)
			if err != nil {
				return err
			}
			fmt.Println(string(data))
			return nil
		}

		if len(processes) == 0 {
			fmt.Println("No processes found.")
			return nil
		}

		for _, p := range processes {
			statusIcon := getStatusIcon(p.Status)
			priorityIcon := getPriorityIcon(p.Priority)
			fmt.Printf("%s [%s] %d: %s\n", statusIcon, priorityIcon, p.ID, p.Title)
			if p.Description != "" {
				fmt.Printf("    └─ %s\n", p.Description)
			}
		}
		return nil
	},
}

func getStatusIcon(status core.ProcessStatus) string {
	switch status {
	case core.StatusRunning:
		return "▶"
	case core.StatusBlocked:
		return "⏸"
	case core.StatusSuspended:
		return "⏹"
	case core.StatusTerminated:
		return "✓"
	default:
		return "?"
	}
}

func getPriorityIcon(priority core.ProcessPriority) string {
	switch priority {
	case core.PriorityHigh:
		return "H"
	case core.PriorityMedium:
		return "M"
	case core.PriorityLow:
		return "L"
	default:
		return "?"
	}
}

func init() {
	rootCmd.AddCommand(psCmd)
	psCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	psCmd.Flags().StringVarP(&psStatus, "status", "s", "running", "Filter by status (use 'all' to show everything)")
	psCmd.Flags().BoolVarP(&psTree, "tree", "t", false, "Display as hierarchical tree")
}
