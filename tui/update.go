package tui

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"taskctl/core"
)

var (
	// Styles
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	cursorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	borderStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))

	// Status styles
	statusStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))  // running - green
	blockedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("226")) // blocked - yellow
	suspendedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241")) // suspended - gray
	doneStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("114")) // terminated - green

	// Priority styles
	priorityHighStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196")) // red
	priorityMediumStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("228")) // yellow
	priorityLowStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("246")) // blue

	// Log styles
	logProgressStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("228")) // yellow
	logStateChangeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))  // green
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case TickMsg:
		// Auto-refresh every tick
		return m, refreshProcesses()

	case ProcessesLoadedMsg:
		m.processes = msg.Processes
		return m, nil

	case ShowDetailMsg:
		return m, m.showProcessDetail(msg.ProcessID)

	case BackToListMsg:
		m.viewMode = ViewList
		m.currentProcess = nil
		m.processLogs = nil
		return m, nil

	case errMsg:
		m.err = msg
		return m, nil

	case ProcessDetailLoadedMsg:
		m.viewMode = ViewDetail
		m.currentProcess = msg.Process
		m.processLogs = msg.Logs
		return m, nil
	}
	return m, nil
}

func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.viewMode == ViewList {
		return m.handleListKeyMsg(msg)
	}
	return m.handleDetailKeyMsg(msg)
}

func (m Model) handleListKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		m.quitting = true
		return m, tea.Quit

	case "j", "down":
		if m.cursor < len(m.processes)-1 {
			m.cursor++
			if m.cursor >= m.viewportOffset+getViewportHeight() {
				m.viewportOffset++
			}
		}
		return m, nil

	case "k", "up":
		if m.cursor > 0 {
			m.cursor--
			if m.cursor < m.viewportOffset {
				m.viewportOffset--
			}
		}
		return m, nil

	case "enter":
		if len(m.processes) > 0 {
			m.selectedIdx = m.cursor
			return m, func() tea.Msg {
				return ShowDetailMsg{ProcessID: m.processes[m.cursor].ID}
			}
		}
		return m, nil

	case "s":
		// TODO: Show spawn dialog
		return m, nil

	case "1", "2", "3", "4", "5", "6", "7", "8", "9":
		// Quick jump to item N-1
		idx := int(msg.String()[0] - '1')
		if idx < len(m.processes) {
			m.cursor = idx
			m.viewportOffset = 0
		}
		return m, nil

	case "?":
		// TODO: Show help
		return m, nil
	}

	return m, nil
}

func (m Model) handleDetailKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		m.quitting = true
		return m, tea.Quit

	case "esc", "h", "left":
		return m, func() tea.Msg { return BackToListMsg{} }

	case "b":
		// Block process
		if m.currentProcess != nil {
			return m, func() tea.Msg {
				_ = core.ChangeProcessState(m.currentProcess.ID, core.StatusBlocked, "")
				return BackToListMsg{}
			}
		}

	case "w":
		// Wake process
		if m.currentProcess != nil {
			return m, func() tea.Msg {
				_ = core.ChangeProcessState(m.currentProcess.ID, core.StatusRunning, "")
				return BackToListMsg{}
			}
		}

	case "t":
		// Terminate process
		if m.currentProcess != nil {
			return m, func() tea.Msg {
				_ = core.ChangeProcessState(m.currentProcess.ID, core.StatusTerminated, "")
				return BackToListMsg{}
			}
		}

	case "a":
		// Add log
		// TODO: Show log input dialog
		return m, nil
	}

	return m, nil
}

func (m Model) showProcessDetail(processID uint) tea.Cmd {
	return func() tea.Msg {
		process, err := core.GetProcess(processID)
		if err != nil {
			return errMsg{err}
		}

		logs, err := core.GetLogs(processID)
		if err != nil {
			return errMsg{err}
		}

		return ProcessDetailLoadedMsg{
			Process: process,
			Logs:    logs,
		}
	}
}

// ProcessDetailLoadedMsg signals that process detail has been loaded
type ProcessDetailLoadedMsg struct {
	Process *core.Process
	Logs    []core.Log
}

func (m Model) handleProcessDetailLoaded(msg ProcessDetailLoadedMsg) (tea.Model, tea.Cmd) {
	m.viewMode = ViewDetail
	m.currentProcess = msg.Process
	m.processLogs = msg.Logs
	return m, nil
}

func getViewportHeight() int {
	return 15 // Approximate visible items
}
