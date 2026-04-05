package tui

import (
	"fmt"

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
	statusStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))  // running - green
	blockedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("226")) // blocked - yellow
	suspendedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241")) // suspended - gray
	doneStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("241")) // terminated - gray

	// Priority styles
	priorityHighStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196")) // red
	priorityMediumStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("228")) // yellow
	priorityLowStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("246")) // blue

	// Log styles
	logProgressStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("228")) // yellow
	logStateChangeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))  // green

	// Warning style
	warningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true) // red bold

	// Markdown styles
	mdBoldStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("228")).Bold(true)       // bold - yellow
	mdItalicStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("228")).Underline(true)  // italic - yellow underline
	mdCodeStyle      = lipgloss.NewStyle().Background(lipgloss.Color("235")).Foreground(lipgloss.Color("252")) // inline code
	mdHeader1Style   = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)        // H1 - green
	mdHeader2Style   = lipgloss.NewStyle().Foreground(lipgloss.Color("80")).Bold(true)        // H2 - cyan
	mdHeader3Style   = lipgloss.NewStyle().Foreground(lipgloss.Color("33")).Bold(true)        // H3 - blue
	mdQuoteStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))                   // quote - gray
	mdLinkStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("117")).Underline(true)  // link - cyan underline
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
		// Return to list view after loading processes
		m.viewMode = ViewList
		return m, nil

	case ProcessDeletedMsg:
		m.processes = msg.Processes
		// Adjust cursor after deletion
		if len(m.processes) == 0 {
			m.cursor = 0
		} else if msg.DeletedIndex >= len(m.processes) {
			// Deleted last item, move to new last item
			m.cursor = len(m.processes) - 1
		} else {
			// Keep cursor at same position (now points to next item)
			m.cursor = msg.DeletedIndex
		}
		// Ensure cursor is within bounds
		if m.cursor >= len(m.processes) {
			m.cursor = len(m.processes) - 1
		}
		m.viewportOffset = 0
		m.viewMode = ViewList
		return m, nil

	case ShowDetailMsg:
		return m, m.showProcessDetail(msg.ProcessID)

	case BackToListMsg:
		m.viewMode = ViewList
		m.currentProcess = nil
		m.processLogs = nil
		// Refresh process list when returning to list view
		return m, refreshProcesses()

	case errMsg:
		m.err = msg
		return m, nil

	case ProcessDetailLoadedMsg:
		m.viewMode = ViewDetail
		m.currentProcess = msg.Process
		m.processLogs = msg.Logs
		m.logCursor = 0 // Reset log cursor
		m.editingLogID = 0 // Reset editing state
		m.pendingStateChange = nil // Clear pending state change
		return m, nil

	case SearchResultsLoadedMsg:
		m.searchKeyword = msg.Keyword
		m.searchResults = msg.Results
		m.searchCursor = 0
		m.viewMode = ViewSearch
		return m, nil

	case TimelineLoadedMsg:
		m.timelineEntries = msg.Entries
		m.timelineCursor = 0
		m.viewMode = ViewTimeline
		return m, nil

	case StatsLoadedMsg:
		m.statsData = msg.Stats
		m.statsDays = msg.Days
		m.viewMode = ViewStats
		return m, nil

	case TreeLoadedMsg:
		m.treeNodes = msg.Nodes
		m.treeCursor = 0
		m.viewMode = ViewTree
		return m, nil

	case ParentsLoadedMsg:
		return m.handleParentsLoaded(msg)

	case ExportSuccessMsg:
		m.exportSuccessMsg = fmt.Sprintf("Exported to: %s", msg.FilePath)
		m.exportFilePath = msg.FilePath
		// Clear export state but keep success message for display
		m.exportProcessID = 0
		m.exportFileName = ""
		// Return to detail view
		m.viewMode = ViewDetail
		return m, nil
	}
	return m, nil
}

func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Debug: record last pressed key (excluding basic typing)
	if msg.Type != tea.KeyRunes {
		m.lastKey = msg.String()
	}
	switch m.viewMode {
	case ViewList:
		return m.handleListKeyMsg(msg)
	case ViewDetail:
		return m.handleDetailKeyMsg(msg)
	case ViewInput:
		return m.handleInputKeyMsg(msg)
	case ViewHelp:
		return m.handleHelpKeyMsg(msg)
	case ViewSpawn, ViewEditProcess:
		return m.handleSpawnKeyMsg(msg)
	case ViewSearch:
		return m.handleSearchKeyMsg(msg)
	case ViewTimeline:
		return m.handleTimelineKeyMsg(msg)
	case ViewStats:
		return m.handleStatsKeyMsg(msg)
	case ViewTree:
		return m.handleTreeKeyMsg(msg)
	case ViewParentSelect:
		return m.handleParentSelectKeyMsg(msg)
	case ViewDeleteConfirm:
		return m.handleDeleteConfirmKeyMsg(msg)
	case ViewExportConfirm:
		return m.handleExportConfirmKeyMsg(msg)
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
		m.editingProcessID = 0 // Ensure we're in create mode
		m.selectedParentID = nil // Reset parent selection
		m.selectedParentName = "" // Reset parent name
		m.viewMode = ViewSpawn
		return m, nil

	case "1", "2", "3", "4", "5", "6", "7", "8", "9":
		// Quick jump to item N-1
		if len(m.processes) > 0 {
			idx := int(msg.String()[0] - '1')
			if idx < len(m.processes) {
				m.cursor = idx
				m.viewportOffset = 0
			}
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

	case "/":
		// Search
		m.textInput.Reset()
		m.textInput.Placeholder = "Search keyword..."
		m.inputPrompt = "Search"
		m.viewMode = ViewInput
		m.editingLogID = 0 // Use editingLogID=0 to indicate search mode
		return m, nil

	case "G":
		// Global timeline
		return m, loadTimeline(0)

	case "D":
		// Stats dashboard
		return m, loadStats(m.statsDays)

	case "Y":
		// Process tree
		return m, loadTree()

	case "d":
		// Delete process with confirmation
		if len(m.processes) > 0 {
			process := m.processes[m.cursor]
			m.confirmDeleteType = "process"
			m.confirmDeleteID = process.ID
			m.confirmDeleteName = process.Title
			m.confirmDeleteIndex = m.cursor // Save cursor position for adjustment after deletion
			m.viewMode = ViewDeleteConfirm
		}
		return m, nil

	case "m":
		// Toggle markdown rendering
		m.markdownEnabled = !m.markdownEnabled
		return m, nil
	}

	return m, nil
}

func (m Model) handleDetailKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		// Return to list instead of quitting
		return m, func() tea.Msg { return BackToListMsg{} }

	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit

	case "esc", "h", "left":
		return m, func() tea.Msg { return BackToListMsg{} }

	case "b":
		// Block process
		if m.currentProcess != nil {
			// Check if already blocked
			if m.currentProcess.Status == core.StatusBlocked {
				return m, nil // Already in this state, do nothing
			}
			// Prompt for optional note
			m.textInput.Reset()
			m.textInput.Placeholder = "输入备注（可选，直接enter跳过）..."
			m.inputPrompt = "阻塞进程"
			status := core.StatusBlocked
			m.pendingStateChange = &status
			m.viewMode = ViewInput
			return m, nil
		}

	case "p":
		// Set process to waiting
		if m.currentProcess != nil {
			// Check if already suspended
			if m.currentProcess.Status == core.StatusSuspended {
				return m, nil // Already in this state, do nothing
			}
			// Prompt for optional note
			m.textInput.Reset()
			m.textInput.Placeholder = "输入备注（可选，直接enter跳过）..."
			m.inputPrompt = "等待中"
			status := core.StatusSuspended
			m.pendingStateChange = &status
			m.viewMode = ViewInput
			return m, nil
		}

	case "w":
		// Wake process
		if m.currentProcess != nil {
			// Check if already running
			if m.currentProcess.Status == core.StatusRunning {
				return m, nil // Already in this state, do nothing
			}
			// Prompt for optional note
			m.textInput.Reset()
			m.textInput.Placeholder = "输入备注（可选，直接enter跳过）..."
			m.inputPrompt = "唤醒进程"
			status := core.StatusRunning
			m.pendingStateChange = &status
			m.viewMode = ViewInput
			return m, nil
		}

	case "t":
		// Terminate process
		if m.currentProcess != nil {
			// Check if already terminated
			if m.currentProcess.Status == core.StatusTerminated {
				return m, nil // Already in this state, do nothing
			}
			// Prompt for optional note
			m.textInput.Reset()
			m.textInput.Placeholder = "输入备注（可选，直接enter跳过）..."
			m.inputPrompt = "终止进程"
			status := core.StatusTerminated
			m.pendingStateChange = &status
			m.viewMode = ViewInput
			return m, nil
		}

	case "a":
		// Add log - show input dialog
		m.textInput.Reset()
		m.textInput.Placeholder = "Enter log content..."
		m.inputPrompt = "Add Log"
		m.editingLogID = 0 // New log
		m.viewMode = ViewInput
		return m, nil

	case "e":
		// Edit selected log
		if len(m.processLogs) > 0 && m.logCursor < len(m.processLogs) {
			log := m.processLogs[m.logCursor]
			m.textInput.Reset()
			m.textInput.SetValue(log.Content)
			m.textInput.Placeholder = "Edit log content..."
			m.inputPrompt = "Edit Log"
			m.editingLogID = log.ID
			m.viewMode = ViewInput
		}
		return m, nil

	case "E":
		// Edit process info
		if m.currentProcess != nil {
			m.spawnTitle.Reset()
			m.spawnTitle.SetValue(m.currentProcess.Title)
			m.spawnDesc.Reset()
			m.spawnDesc.SetValue(m.currentProcess.Description)
			m.spawnPriority.Reset()
			m.spawnPriority.SetValue(string(m.currentProcess.Priority))
			m.spawnFocusedField = 0
			m.spawnTitle.Focus()
			m.editingProcessID = m.currentProcess.ID
			m.selectedParentID = m.currentProcess.ParentID // Set current parent
			// Cache parent name for display
			m.selectedParentName = ""
			m.availableParents = []core.Process{} // Clear stale parents
			m.viewMode = ViewEditProcess
		}
		return m, nil

	case "x":
		// Delete selected log with confirmation
		if len(m.processLogs) > 0 && m.logCursor < len(m.processLogs) {
			log := m.processLogs[m.logCursor]
			m.confirmDeleteType = "log"
			m.confirmDeleteID = log.ID
			// Truncate log content for display
			logContent := log.Content
			if len(logContent) > 30 {
				logContent = logContent[:30] + "..."
			}
			m.confirmDeleteName = logContent
			m.viewMode = ViewDeleteConfirm
		}
		return m, nil

	case "j", "down":
		// Move cursor down in logs
		if m.logCursor < len(m.processLogs)-1 {
			m.logCursor++
		}
		return m, nil

	case "k", "up":
		// Move cursor up in logs
		if m.logCursor > 0 {
			m.logCursor--
		}
		return m, nil

	case "m":
		// Toggle markdown rendering
		m.markdownEnabled = !m.markdownEnabled
		return m, nil

	case "c":
		// Copy log content to clipboard (placeholder for future implementation)
		// TODO: Implement clipboard functionality using atotto/clipboard
		return m, nil

	case ">":
		// Export process as Markdown
		if m.currentProcess != nil {
			m.exportProcessID = m.currentProcess.ID
			m.exportFileName = GenerateExportFileName(m.currentProcess)
			m.exportFilePath = GetAbsolutePath(m.exportFileName)
			m.viewMode = ViewExportConfirm
		}
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
	// Debug: record key for input view
	m.lastKey = msg.String()

	switch msg.String() {
	case "ctrl+enter", "ctrl+j":
		// Ctrl+Enter or Ctrl+J (macOS compatible) to submit log
		if m.currentProcess != nil && m.pendingStateChange == nil {
			input := m.textInput.Value()
			if input == "" {
				return m, nil
			}
			return m, func() tea.Msg {
				var err error

				// Check if editing existing log or adding new one
				if m.editingLogID > 0 {
					// Update existing log
					err = core.UpdateLog(m.editingLogID, input)
				} else {
					// Add new log
					_, err = core.AddLog(m.currentProcess.ID, core.LogTypeProgress, input)
				}

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

	case "enter":
		// For state change: Enter submits (note is optional)
		if m.pendingStateChange != nil && m.currentProcess != nil {
			note := m.textInput.Value()
			newStatus := *m.pendingStateChange
			return m, func() tea.Msg {
				if err := core.ChangeProcessState(m.currentProcess.ID, newStatus, note); err != nil {
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

		// For search: Enter submits
		if m.inputPrompt == "Search" {
			input := m.textInput.Value()
			if input == "" {
				return m, nil
			}
			return m, loadSearchResults(input)
		}

		// For log input: let textarea handle Enter as newline
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd

	case "esc":
		// Cancel input - return to previous view
		if m.inputPrompt == "Search" {
			m.viewMode = ViewList
			return m, refreshProcesses()
		} else {
			m.viewMode = ViewDetail
		}
		m.editingLogID = 0
		m.pendingStateChange = nil // Clear pending state change
		return m, nil

	default:
		// Update textarea
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}
}

func (m Model) handleHelpKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc", "?":
		m.viewMode = ViewList
		m.helpOffset = 0 // Reset scroll
		return m, nil

	case "j", "down":
		m.helpOffset++
		return m, nil

	case "k", "up":
		if m.helpOffset > 0 {
			m.helpOffset--
		}
		return m, nil
	}
	return m, nil
}

func (m Model) handleSpawnKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Debug: record key for spawn view
	m.lastKey = msg.String()

	switch msg.String() {
	case "esc":
		// Cancel - return to appropriate view
		if m.editingProcessID > 0 {
			// Was editing, return to detail view
			m.viewMode = ViewDetail
			return m, func() tea.Msg {
				process, err := core.GetProcess(m.editingProcessID)
				if err != nil {
					return errMsg{err}
				}
				logs, err := core.GetLogs(m.editingProcessID)
				if err != nil {
					return errMsg{err}
				}
				return ProcessDetailLoadedMsg{Process: process, Logs: logs}
			}
		}
		// Was creating new, return to list view
		m.viewMode = ViewList
		return m, refreshProcesses()

	case "ctrl+enter", "ctrl+j":
		// Ctrl+Enter or Ctrl+J (macOS compatible) submits the form
		return m.submitSpawnForm()

	case "enter":
		// Regular Enter on parent field - open parent selection
		if m.spawnFocusedField == 3 {
			return m, m.loadAvailableParents()
		}

		// For description field (textarea), let it handle enter (newline)
		if m.spawnFocusedField == 1 {
			// Update the textarea with the enter key
			var cmd tea.Cmd
			m.spawnDesc, cmd = m.spawnDesc.Update(msg)
			return m, cmd
		}

		// For other fields, no action on enter (need Ctrl+Enter to submit)
		return m, nil

	case "tab":
		// Navigate between fields
		m.spawnFocusedField = (m.spawnFocusedField + 1) % 4
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
		case 3:
			m.spawnTitle.Blur()
			m.spawnDesc.Blur()
			m.spawnPriority.Blur()
		}
		return m, nil

	case "shift+tab":
		// Navigate backwards
		m.spawnFocusedField = (m.spawnFocusedField + 3) % 4
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
		case 3:
			m.spawnTitle.Blur()
			m.spawnDesc.Blur()
			m.spawnPriority.Blur()
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
		case 3:
			// Parent field is read-only (no text input)
		}
		return m, cmd
	}
}

// submitSpawnForm submits the spawn/edit form
func (m Model) submitSpawnForm() (tea.Model, tea.Cmd) {
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

	if m.editingProcessID > 0 {
		// Update existing process
		return m, func() tea.Cmd {
			err := core.UpdateProcess(m.editingProcessID, &title, &desc, &priority, m.selectedParentID)
			if err != nil {
				return func() tea.Msg { return errMsg{err} }
			}
			// Refresh process detail
			process, err := core.GetProcess(m.editingProcessID)
			if err != nil {
				return func() tea.Msg { return errMsg{err} }
			}
			logs, err := core.GetLogs(m.editingProcessID)
			if err != nil {
				return func() tea.Msg { return errMsg{err} }
			}
			// Also refresh the process list to show updated info
			return tea.Batch(
				func() tea.Msg { return ProcessDetailLoadedMsg{Process: process, Logs: logs} },
				refreshProcessesWithFilter(m.statusFilter),
			)
		}()
	}

	// Create new process
	return m, func() tea.Msg {
		_, err := core.CreateProcess(title, desc, m.selectedParentID, priority)
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
}

func (m Model) handleSearchKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		// Return to list
		m.viewMode = ViewList
		return m, refreshProcesses()

	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit

	case "j", "down":
		if m.searchCursor < len(m.searchResults)-1 {
			m.searchCursor++
		}
		return m, nil

	case "k", "up":
		if m.searchCursor > 0 {
			m.searchCursor--
		}
		return m, nil

	case "enter":
		if len(m.searchResults) > 0 && m.searchCursor < len(m.searchResults) {
			result := m.searchResults[m.searchCursor]
			if result.Type == "process" {
				return m, func() tea.Msg {
					return ShowDetailMsg{ProcessID: result.ID}
				}
			} else {
				// It's a log, navigate to its process
				return m, func() tea.Msg {
					return ShowDetailMsg{ProcessID: result.ProcessID}
				}
			}
		}
		return m, nil

	case "/":
		// New search
		m.textInput.Reset()
		m.textInput.Placeholder = "Search keyword..."
		m.inputPrompt = "Search"
		m.viewMode = ViewInput
		return m, nil
	}
	return m, nil
}

func (m Model) handleTimelineKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		// Return to list
		m.viewMode = ViewList
		return m, refreshProcesses()

	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit

	case "j", "down":
		if m.timelineCursor < len(m.timelineEntries)-1 {
			m.timelineCursor++
		}
		return m, nil

	case "k", "up":
		if m.timelineCursor > 0 {
			m.timelineCursor--
		}
		return m, nil

	case "enter":
		if len(m.timelineEntries) > 0 && m.timelineCursor < len(m.timelineEntries) {
			entry := m.timelineEntries[m.timelineCursor]
			return m, func() tea.Msg {
				return ShowDetailMsg{ProcessID: entry.ProcessID}
			}
		}
		return m, nil
	}
	return m, nil
}

func (m Model) handleStatsKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		// Return to list
		m.viewMode = ViewList
		return m, refreshProcesses()

	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit

	case "h":
		// Show stats for past 7 days
		return m, loadStats(7)

	case "H":
		// Show stats for past 90 days
		return m, loadStats(90)

	case "d":
		// Show stats for past 30 days (default)
		return m, loadStats(30)

	case "D":
		// Show stats for past 365 days
		return m, loadStats(365)
	}
	return m, nil
}

func (m Model) handleTreeKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		// Return to list
		m.viewMode = ViewList
		return m, refreshProcesses()

	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit

	case "j", "down":
		// Navigate to next visible node
		m.treeCursor = m.getNextVisibleTreeNode(m.treeCursor, 1)
		return m, nil

	case "k", "up":
		// Navigate to previous visible node
		m.treeCursor = m.getNextVisibleTreeNode(m.treeCursor, -1)
		return m, nil

	case "enter":
		// View process details
		nodes := m.getVisibleTreeNodes()
		if m.treeCursor < len(nodes) {
			return m, func() tea.Msg {
				return ShowDetailMsg{ProcessID: nodes[m.treeCursor].ID}
			}
		}
		return m, nil

	case " ":
		// Toggle expand/collapse
		nodes := m.getVisibleTreeNodes()
		if m.treeCursor < len(nodes) {
			nodeID := nodes[m.treeCursor].ID
			if m.treeExpanded[nodeID] {
				delete(m.treeExpanded, nodeID)
			} else {
				m.treeExpanded[nodeID] = true
			}
		}
		return m, nil
	}
	return m, nil
}

// getVisibleTreeNodes returns all visible tree nodes (respecting expand/collapse)
func (m Model) getVisibleTreeNodes() []core.ProcessNode {
	return m.flattenTreeNodes(m.treeNodes, 0)
}

// flattenTreeNodes recursively flattens tree nodes to a list
func (m Model) flattenTreeNodes(nodes []*core.ProcessNode, depth int) []core.ProcessNode {
	var result []core.ProcessNode
	for _, node := range nodes {
		// Add the node
		result = append(result, *node)
		// Add children if expanded
		if m.treeExpanded[node.ID] && len(node.Children) > 0 {
			result = append(result, m.flattenTreeNodes(node.Children, depth+1)...)
		}
	}
	return result
}

// getNextVisibleTreeNode finds the next/previous visible tree node
func (m Model) getNextVisibleTreeNode(current int, direction int) int {
	nodes := m.getVisibleTreeNodes()
	if len(nodes) == 0 {
		return 0
	}

	newIndex := current + direction
	if newIndex < 0 {
		newIndex = 0
	} else if newIndex >= len(nodes) {
		newIndex = len(nodes) - 1
	}

	return newIndex
}

// loadAvailableParents loads the available parent processes and switches to parent selection view
func (m Model) loadAvailableParents() tea.Cmd {
	return func() tea.Msg {
		parents, err := core.ListProcesses(nil)
		if err != nil {
			return errMsg{err}
		}
		return ParentsLoadedMsg{Parents: parents}
	}
}

// handleParentsLoaded handles the ParentsLoadedMsg
func (m Model) handleParentsLoaded(msg ParentsLoadedMsg) (tea.Model, tea.Cmd) {
	// Filter out the current process if editing (can't be parent of itself)
	if m.editingProcessID > 0 {
		// Get all descendant IDs to filter out
		descendantIDs, err := core.GetDescendantIDs(m.editingProcessID)
		if err != nil {
			// If error, just filter self
			descendantIDs = []uint{m.editingProcessID}
		}
		// Also add self to the filter list
		descendantIDs = append(descendantIDs, m.editingProcessID)

		// Create a set for faster lookup
		filterSet := make(map[uint]bool)
		for _, id := range descendantIDs {
			filterSet[id] = true
		}

		filtered := make([]core.Process, 0, len(msg.Parents))
		for _, p := range msg.Parents {
			if !filterSet[p.ID] {
				filtered = append(filtered, p)
			}
			// Cache parent name if this is the selected parent
			if m.selectedParentID != nil && p.ID == *m.selectedParentID {
				m.selectedParentName = p.Title
			}
		}
		m.availableParents = filtered
	} else {
		m.availableParents = msg.Parents
		// Cache parent name for new process
		if m.selectedParentID != nil {
			for _, p := range msg.Parents {
				if p.ID == *m.selectedParentID {
					m.selectedParentName = p.Title
					break
				}
			}
		}
	}

	m.parentCursor = -1 // Start at "None" option
	m.viewMode = ViewParentSelect
	return m, nil
}

func (m Model) handleParentSelectKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		// Return to spawn/edit view
		if m.editingProcessID > 0 {
			m.viewMode = ViewEditProcess
		} else {
			m.viewMode = ViewSpawn
		}
		return m, nil

	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit

	case "j", "down":
		// Move cursor down
		maxIdx := len(m.availableParents) // -1 to "None" is at cursor=-1
		if m.parentCursor < maxIdx {
			m.parentCursor++
		}
		return m, nil

	case "k", "up":
		// Move cursor up
		if m.parentCursor > -1 {
			m.parentCursor--
		}
		return m, nil

	case "enter":
		// Select the parent and return to spawn/edit view
		if m.parentCursor == -1 {
			// Selected "None"
			m.selectedParentID = nil
			m.selectedParentName = ""
		} else if m.parentCursor < len(m.availableParents) {
			// Selected a parent
			parentID := m.availableParents[m.parentCursor].ID
			m.selectedParentID = &parentID
			m.selectedParentName = m.availableParents[m.parentCursor].Title
		}

		// Return to spawn/edit view
		if m.editingProcessID > 0 {
			m.viewMode = ViewEditProcess
		} else {
			m.viewMode = ViewSpawn
		}
		return m, nil
	}
	return m, nil
}

func (m Model) handleDeleteConfirmKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		// Confirm deletion
		if m.confirmDeleteType == "process" {
			// Clear confirmation state and switch to list view immediately
			deleteID := m.confirmDeleteID
			deletedIndex := m.confirmDeleteIndex
			currentFilter := m.statusFilter
			m.confirmDeleteType = ""
			m.confirmDeleteID = 0
			m.confirmDeleteName = ""
			m.confirmDeleteIndex = 0
			m.viewMode = ViewList
			return m, func() tea.Msg {
				err := core.DeleteProcess(deleteID)
				if err != nil {
					return errMsg{err}
				}
				// Refresh processes with current filter
				var processes []core.Process
				if currentFilter == "" {
					processes, err = core.ListProcesses(nil)
				} else {
					processes, err = core.ListProcesses(&currentFilter)
				}
				if err != nil {
					return errMsg{err}
				}
				return ProcessDeletedMsg{Processes: processes, DeletedIndex: deletedIndex}
			}
		} else if m.confirmDeleteType == "log" && m.currentProcess != nil {
			// Clear confirmation state and switch to detail view immediately
			deleteID := m.confirmDeleteID
			processID := m.currentProcess.ID
			m.confirmDeleteType = ""
			m.confirmDeleteID = 0
			m.confirmDeleteName = ""
			m.viewMode = ViewDetail
			return m, func() tea.Msg {
				err := core.DeleteLog(deleteID)
				if err != nil {
					return errMsg{err}
				}
				// Refresh process detail
				logs, err := core.GetLogs(processID)
				if err != nil {
					return errMsg{err}
				}
				return ProcessDetailLoadedMsg{
					Process: m.currentProcess,
					Logs:    logs,
				}
			}
		}
		return m, nil

	case "n", "N", "q", "esc":
		// Cancel deletion - return to previous view
		if m.confirmDeleteType == "process" {
			m.viewMode = ViewList
			// Clear confirmation state
			m.confirmDeleteType = ""
			m.confirmDeleteID = 0
			m.confirmDeleteName = ""
			return m, refreshProcesses()
		} else {
			m.viewMode = ViewDetail
		}
		// Clear confirmation state
		m.confirmDeleteType = ""
		m.confirmDeleteID = 0
		m.confirmDeleteName = ""
		return m, nil

	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	}
	return m, nil
}

// handleExportConfirmKeyMsg handles key messages in export confirmation view
func (m Model) handleExportConfirmKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		// Confirm export
		exportID := m.exportProcessID
		m.viewMode = ViewDetail
		return m, func() tea.Msg {
			filePath, err := ExportProcess(exportID)
			if err != nil {
				return errMsg{err}
			}
			return ExportSuccessMsg{FilePath: filePath}
		}

	case "n", "N", "q", "esc":
		// Cancel export - return to detail view
		m.viewMode = ViewDetail
		// Clear export state
		m.exportProcessID = 0
		m.exportFileName = ""
		m.exportFilePath = ""
		return m, nil

	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	}
	return m, nil
}
