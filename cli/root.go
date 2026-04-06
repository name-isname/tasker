package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"taskctl/core"
	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
	localDB    bool

	// Version information injected by goreleaser
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

// Command groups for better organization
const (
	GroupProcess = "process"
	GroupState   = "state"
	GroupLogs    = "logs"
	GroupAnalysis = "analysis"
	GroupUI      = "ui"
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

		// Determine database path (priority: --db > --local > default global)
		dbPath, _ := cmd.Flags().GetString("db")
		if dbPath == "./taskctl.db" {
			// --db not explicitly set, check --local flag
			local, _ := cmd.Flags().GetBool("local")
			if !local {
				// Use global path
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return err
				}
				dbPath = filepath.Join(homeDir, ".taskctl", "taskctl.db")
			}
		}

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
func Execute(version, commit, date string) {
	Version = version
	Commit = commit
	BuildDate = date

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("db", "d", "./taskctl.db", "Path to the SQLite database")
	rootCmd.PersistentFlags().BoolVarP(&localDB, "local", "L", false, "Use local database in current directory")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output as JSON")

	// Set up command groups for better organization
	rootCmd.AddGroup(&cobra.Group{
		ID:    GroupProcess,
		Title: "Process Management",
	})

	rootCmd.AddGroup(&cobra.Group{
		ID:    GroupState,
		Title: "State Management",
	})

	rootCmd.AddGroup(&cobra.Group{
		ID:    GroupLogs,
		Title: "Log Management",
	})

	rootCmd.AddGroup(&cobra.Group{
		ID:    GroupAnalysis,
		Title: "Analysis & Search",
	})

	rootCmd.AddGroup(&cobra.Group{
		ID:    GroupUI,
		Title: "Interface & Export",
	})

	// Add version command
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("taskctl version %s\n", Version)
			fmt.Printf("commit: %s\n", Commit)
			fmt.Printf("built at: %s\n", BuildDate)
		},
	}
	rootCmd.AddCommand(versionCmd)

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
