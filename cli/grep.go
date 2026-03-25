package cli

import (
	"fmt"
	"strings"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var grepCmd = &cobra.Command{
	Use:   "grep <keyword>",
	Short: "Search across processes and logs",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		keyword := args[0]

		results, err := core.GlobalSearch(keyword)
		if err != nil {
			return err
		}

		if len(results) == 0 {
			fmt.Printf("No results found for '%s'\n", keyword)
			return nil
		}

		fmt.Printf("\nFound %d results for '%s':\n\n", len(results), keyword)

		for _, result := range results {
			if result.Type == "process" {
				fmt.Printf("[Process #%d] %s\n", result.ID, result.Title)
				if result.Content != "" {
					// Show preview
					preview := result.Content
					if len(preview) > 60 {
						preview = preview[:60] + "..."
					}
					fmt.Printf("    └─ %s\n", preview)
				}
			} else {
				fmt.Printf("[Log #%d] in Process #%d\n", result.ID, result.ProcessID)
				// Highlight matching content
				content := result.Content
				if len(content) > 80 {
					// Try to find keyword and show context
					idx := strings.Index(strings.ToLower(content), strings.ToLower(keyword))
					if idx >= 0 {
						start := idx - 20
						if start < 0 {
							start = 0
						}
						end := idx + len(keyword) + 20
						if end > len(content) {
							end = len(content)
						}
						content = "..." + content[start:end] + "..."
					} else {
						content = content[:80] + "..."
					}
				}
				fmt.Printf("    └─ %s\n", content)
			}
			fmt.Println()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(grepCmd)
	grepCmd.GroupID = GroupAnalysis
}
