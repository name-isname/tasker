package tui

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"taskctl/core"
)

var (
	cursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Faint(true)
	statusStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))  // running
	blockedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("226")) // blocked
	suspendedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")) // suspended
	doneStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("114")) // terminated
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.Quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.Selected > 0 {
				m.Selected--
			}
		case "down", "j":
			if m.Selected < len(m.Processes)-1 {
				m.Selected++
			}
		}
	}
	return m, nil
}

// View renders the TUI
func (m Model) View() string {
	if m.Quitting {
		return "Goodbye!\n"
	}

	s := "Process List\n\n"

	for i, process := range m.Processes {
		cursor := " "
		if m.Selected == i {
			cursor = ">"
		}

		// Style based on status
		var statusStyle lipgloss.Style
		var statusIcon string
		switch process.Status {
		case core.StatusRunning:
			statusStyle = statusStyle
			statusIcon = "▶"
		case core.StatusBlocked:
			statusStyle = blockedStyle
			statusIcon = "⏸"
		case core.StatusSuspended:
			statusStyle = suspendedStyle
			statusIcon = "⏹"
		case core.StatusTerminated:
			statusStyle = doneStyle
			statusIcon = "✓"
		default:
			statusStyle = helpStyle
			statusIcon = "?"
		}

		line := cursorStyle.Render(cursor) + " " +
			statusStyle.Render(statusIcon) + " " +
			process.Title
		s += line + "\n"

		if m.Selected == i && process.Description != "" {
			s += helpStyle.Render("    └─ "+process.Description) + "\n"
		}
	}

	s += "\n" + helpStyle.Render("Press q to quit")
	return s
}
