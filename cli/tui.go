package cli

import (
	"os"
	"taskctl/tui"
	"github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the terminal UI",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(tui.InitialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
		_, err := p.Run()
		if err != nil {
			return err
		}
		os.Exit(0)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
	tuiCmd.GroupID = GroupUI
}
