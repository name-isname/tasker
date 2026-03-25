package tui

import (
	"github.com/charmbracelet/bubbletea"
	"taskctl/core"
)

// Model represents the TUI state
type Model struct {
	Processes []core.Process
	Selected  int
	Quitting  bool
}

// Init initializes the TUI model
func (m Model) Init() tea.Cmd {
	return nil
}
