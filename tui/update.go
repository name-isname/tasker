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
		// Auto-refresh every tick with current filter
		return m, refreshProcessesWithFilter(m.statusFilter)

	case ProcessesLoadedMsg:
		m.processes = msg.Processes
		// Reset cursor if out of bounds
		if m.cursor >= len(m.processes) {
			m.cursor = 0
			m.viewportOffset = 0
		}
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
	switch m.viewMode {
	case ViewList:
		return m.handleListKeyMsg(msg)
	case ViewDetail:
		return m.handleDetailKeyMsg(msg)
	case ViewInput:
		return m.handleInputKeyMsg(msg)
	case ViewHelp:
		return m.handleHelpKeyMsg(msg)
	case ViewSpawn:
		return m.handleSpawnKeyMsg(msg)
	}
	return m, nil
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
		// Show spawn dialog
		m.spawnTitle.Reset()
		m.spawnDesc.Reset()
		m.spawnPriority.SetValue("M")
		m.spawnFocusedField = 0
		m.spawnTitle.Focus()
		m.viewMode = ViewSpawn
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
		m.viewMode = ViewHelp
		return m, nil

	case "r":
		// Filter by running
		m.statusFilter = core.StatusRunning
		return m, refreshProcessesWithFilter(core.StatusRunning)

	case "B":
		// Filter by blocked
		m.statusFilter = core.StatusBlocked
		return m, refreshProcessesWithFilter(core.StatusBlocked)

	case "S":
		// Filter by suspended
		m.statusFilter = core.StatusSuspended
		return m, refreshProcessesWithFilter(core.StatusSuspended)

	case "T":
		// Filter by terminated
		m.statusFilter = core.StatusTerminated
		return m, refreshProcessesWithFilter(core.StatusTerminated)

	case "A":
		// Show all
		m.statusFilter = ""
		return m, refreshProcessesWithFilter("")
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
		// Add log - show input dialog
		m.textInput.Reset()
		m.textInput.Placeholder = "Enter log content..."
		m.inputPrompt = "Add Log"
		m.viewMode = ViewInput
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

func (m Model) handleInputKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Submit the input
		input := m.textInput.Value()
		if input != "" && m.currentProcess != nil {
			return m, func() tea.Msg {
				_, err := core.AddLog(m.currentProcess.ID, core.LogTypeProgress, input)
				if err != nil {
					return errMsg{err}
				}
				// Refresh process detail
				process, err := core.GetProcess(m.currentProcess.ID)
				if err != nil {
					return errMsg{err}
				}
				logs, err := core.GetLogs(m.currentProcess.ID)
				if err != nil {
					return errMsg{err}
				}
				return ProcessDetailLoadedMsg{
					Process: process,
					Logs:    logs,
				}
			}
		}
		return m, nil

	case "esc":
		// Cancel input
		m.viewMode = ViewDetail
		return m, nil

	default:
		// Update text input
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}
}

func (m Model) handleHelpKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc", "?":
		m.viewMode = ViewList
		return m, nil
	}
	return m, nil
}

func (m Model) handleSpawnKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		// Cancel
		m.viewMode = ViewList
		return m, nil

	case "enter":
		// Submit and create process
		title := m.spawnTitle.Value()
		desc := m.spawnDesc.Value()
		priorityStr := m.spawnPriority.Value()

		if title == "" {
			return m, nil
		}

		var priority core.ProcessPriority
		switch priorityStr {
		case "H", "h":
			priority = core.PriorityHigh
		case "L", "l":
			priority = core.PriorityLow
		default:
			priority = core.PriorityMedium
		}

		return m, func() tea.Msg {
			_, err := core.CreateProcess(title, desc, nil, priority)
			if err != nil {
				return errMsg{err}
			}
			// Refresh processes
			processes, err := core.ListProcesses(nil)
			if err != nil {
				return errMsg{err}
			}
			return ProcessesLoadedMsg{Processes: processes}
		}

	case "tab":
		// Navigate between fields
		m.spawnFocusedField = (m.spawnFocusedField + 1) % 3
		switch m.spawnFocusedField {
		case 0:
			m.spawnTitle.Focus()
			m.spawnDesc.Blur()
			m.spawnPriority.Blur()
		case 1:
			m.spawnTitle.Blur()
			m.spawnDesc.Focus()
			m.spawnPriority.Blur()
		case 2:
			m.spawnTitle.Blur()
			m.spawnDesc.Blur()
			m.spawnPriority.Focus()
		}
		return m, nil

	case "shift+tab":
		// Navigate backwards
		m.spawnFocusedField = (m.spawnFocusedField + 2) % 3
		switch m.spawnFocusedField {
		case 0:
			m.spawnTitle.Focus()
			m.spawnDesc.Blur()
			m.spawnPriority.Blur()
		case 1:
			m.spawnTitle.Blur()
			m.spawnDesc.Focus()
			m.spawnPriority.Blur()
		case 2:
			m.spawnTitle.Blur()
			m.spawnDesc.Blur()
			m.spawnPriority.Focus()
		}
		return m, nil

	default:
		// Update the focused field
		var cmd tea.Cmd
		switch m.spawnFocusedField {
		case 0:
			m.spawnTitle, cmd = m.spawnTitle.Update(msg)
		case 1:
			m.spawnDesc, cmd = m.spawnDesc.Update(msg)
		case 2:
			m.spawnPriority, cmd = m.spawnPriority.Update(msg)
		}
		return m, cmd
	}
}
