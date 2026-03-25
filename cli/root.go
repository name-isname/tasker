package cli

import (
	"os"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "taskctl",
	Short: "A task management tool with CLI, TUI, and Web interfaces",
	Long:  `taskctl is a 3-in-1 task management tool supporting CLI, TUI, and Web UI interfaces.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize database before running commands
		dbPath, _ := cmd.Flags().GetString("db")
		if err := core.InitDB(dbPath); err != nil {
			return err
		}
		return core.AutoMigrate()
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("db", "./taskctl.db", "Path to the SQLite database")
}
