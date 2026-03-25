package cli

import (
	"encoding/json"
	"fmt"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var jsonOutput bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		tasks, err := core.ListTasks()
		if err != nil {
			return err
		}

		if jsonOutput {
			data, err := json.Marshal(tasks)
			if err != nil {
				return err
			}
			fmt.Println(string(data))
			return nil
		}

		for _, task := range tasks {
			status := " "
			if task.Completed {
				status = "x"
			}
			fmt.Printf("[%s] %d: %s\n", status, task.ID, task.Title)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
}
