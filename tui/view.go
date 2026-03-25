package tui

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	cursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Faint(true)
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
			if m.Selected < len(m.Tasks)-1 {
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

	s := "Task List\n\n"

	for i, task := range m.Tasks {
		cursor := " "
		if m.Selected == i {
			cursor = ">"
		}
		status := " "
		if task.Completed {
			status = "x"
		}
		s += cursorStyle.Render(cursor) + " [" + status + "] " + task.Title + "\n"
	}

	s += "\n" + helpStyle.Render("Press q to quit")
	return s
}
