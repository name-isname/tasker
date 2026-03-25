package cli

import (
	"os"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
)

var rootCmd = &cobra.Command{
	Use:   "taskctl",
	Short: "A process-oriented task management tool",
	Long: `taskctl is a 3-in-1 task management tool supporting CLI, TUI, and Web UI interfaces.

It models tasks as OS "Processes" with state transitions (running, blocked, suspended, terminated)
rather than simple todo items. Every state change and progress note is recorded as a chronological Log.`,
	Example: `  # Create a new process
  taskctl spawn "Build web app" -D "Create personal website" -P high

  # List all running processes
  taskctl ps

  # Show process tree
  taskctl ps -t

  # Inspect a process
  taskctl inspect 1

  # Add progress log
  taskctl log 1 "Started implementation"

  # Block a process with reason
  taskctl block 1 -m "Waiting for API key"

  # Wake up a process
  taskctl wake 1 -m "API key received"

  # Search across all processes
  taskctl grep "database"

  # View global timeline
  taskctl timeline

  # Export as Markdown
  taskctl export 1

  # Output as JSON (for AI agents)
  taskctl ps --json`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip DB init for help command
		if cmd.Name() == "help" || len(args) > 0 && args[0] == "help" {
			return nil
		}

		// Initialize database before running commands
		dbPath, _ := cmd.Flags().GetString("db")
		if err := core.InitDB(dbPath); err != nil {
			return err
		}
		return core.AutoMigrate()
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("db", "d", "./taskctl.db", "Path to the SQLite database")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output as JSON")

	// Add custom help flag that works with -h
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:   "help [command]",
		Short: "Print usage information",
		Long:  `Print usage information for any command or subcommand.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return rootCmd.Help()
			}
			c, _, err := rootCmd.Find(args)
			if err != nil {
				return err
			}
			return c.Help()
		},
	})
}
