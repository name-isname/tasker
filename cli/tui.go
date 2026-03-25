package cli

import (
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the terminal UI",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Launch TUI with tea.NewProgram()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
	tuiCmd.GroupID = GroupUI
}
