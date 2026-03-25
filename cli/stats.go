package cli

import (
	"fmt"
	"strings"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var (
	statsDays int
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show activity statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		days := statsDays
		if days == 0 {
			days = 30 // Default to 30 days
		}

		stats, err := core.GetActivityStats(days)
		if err != nil {
			return err
		}

		if len(stats) == 0 {
			fmt.Printf("No activity in the past %d days.\n", days)
			return nil
		}

		// Find max count for scaling
		maxCount := 0
		for _, stat := range stats {
			if stat.Count > maxCount {
				maxCount = stat.Count
			}
		}

		// Print activity heatmap
		fmt.Printf("\n📊 Activity (past %d days)\n", days)
		fmt.Println(strings.Repeat("-", 50))

		for _, stat := range stats {
			// Create a simple bar chart
			barWidth := int(float64(stat.Count) / float64(maxCount) * 30)
			bar := strings.Repeat("█", barWidth)

			// Format date as MM-DD
			dateStr := stat.Date[5:] // Skip YYYY-

			fmt.Printf("%s │%s %d\n", dateStr, bar, stat.Count)
		}

		// Calculate totals
		totalLogs := 0
		for _, stat := range stats {
			totalLogs += stat.Count
		}
		avgPerDay := float64(totalLogs) / float64(len(stats))

		fmt.Println(strings.Repeat("-", 50))
		fmt.Printf("Total: %d logs | Average: %.1f logs/day\n", totalLogs, avgPerDay)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
	statsCmd.GroupID = GroupAnalysis
	statsCmd.Flags().IntVarP(&statsDays, "days", "D", 30, "Number of days to analyze")
}
