package cli

import (
	"fmt"
	"strings"
	"taskctl/core"
	"time"
	"github.com/spf13/cobra"
)

var (
	timelineDays  int
	timelineLimit int
)

var timelineCmd = &cobra.Command{
	Use:   "timeline",
	Short: "Show global timeline of all logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		var entries []core.TimelineEntry
		var err error

		if timelineDays > 0 {
			// Get timeline for past N days
			startTime := time.Now().AddDate(0, 0, -timelineDays)
			entries, err = core.GetTimeline(startTime, time.Time{}, timelineLimit)
		} else {
			// Get today's timeline by default
			entries, err = core.GetTodayTimeline()
			if len(entries) == 0 {
				fmt.Println("No logs for today.")
				return nil
			}
		}

		if err != nil {
			return err
		}

		if len(entries) == 0 {
			fmt.Println("No logs found.")
			return nil
		}

		// Group entries by date
		grouped := make(map[string][]core.TimelineEntry)
		for _, entry := range entries {
			date := entry.CreatedAt.Format("2006-01-02")
			grouped[date] = append(grouped[date], entry)
		}

		// Print timeline grouped by date
		for date, dateEntries := range grouped {
			fmt.Printf("\n📅 %s\n", date)
			fmt.Println(strings.Repeat("-", 40))

			for _, entry := range dateEntries {
				timeStr := entry.CreatedAt.Format("15:04")
				icon := getLogIcon(entry.LogType)

				fmt.Printf("%s [%s] %s: %s\n",
					timeStr, icon, entry.ProcessTitle, entry.Content)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(timelineCmd)
	timelineCmd.Flags().IntVarP(&timelineDays, "days", "d", 0, "Show logs for past N days")
	timelineCmd.Flags().IntVarP(&timelineLimit, "limit", "n", 50, "Limit number of entries")
}

func getLogIcon(logType core.LogType) string {
	switch logType {
	case core.LogTypeStateChange:
		return "🔄"
	case core.LogTypeProgress:
		return "📝"
	default:
		return "📌"
	}
}
