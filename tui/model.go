package tui

import (
	"time"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"taskctl/core"
)

// ViewMode represents the current view state
type ViewMode int

const (
	ViewList ViewMode = iota
	ViewDetail
	ViewInput
	ViewHelp
	ViewSpawn
	ViewEditProcess
)

// Model represents the TUI state
type Model struct {
	// View state
	viewMode    ViewMode
	processes   []core.Process
	selectedIdx int
	quitting    bool

	// Detail view state
	currentProcess *core.Process
	processLogs    []core.Log

	// Viewport and cursor
	viewportOffset int
	cursor         int

	// Error message
	err error

	// Status filter
	statusFilter core.ProcessStatus
	filtering     bool

	// Text input
	textInput     textinput.Model
	inputPrompt   string
	editingLogID  uint // 0 for new log, >0 for editing existing log

	// Spawn/Edit form fields (shared between spawn and edit)
	spawnTitle       textinput.Model
	spawnDesc        textinput.Model
	spawnPriority    textinput.Model
	spawnFocusedField int // 0=title, 1=desc, 2=priority
	editingProcessID  uint // 0 for new process, >0 for editing existing process

	// Log selection in detail view
	logCursor int // Index of selected log in processLogs
}

// Messages
type (
	TickMsg      time.Time
	RefreshMsg   struct{}
	ShowDetailMsg struct {
		ProcessID uint
	}
	BackToListMsg struct{}
	StatusChangeMsg struct {
		ProcessID uint
		Status    core.ProcessStatus
	}
)

// InitialModel creates the initial TUI model
func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter log content..."
	ti.Focus()

	titleInput := textinput.New()
	titleInput.Placeholder = "Process title..."

	descInput := textinput.New()
	descInput.Placeholder = "Description (optional)..."

	priorityInput := textinput.New()
	priorityInput.Placeholder = "M"
	priorityInput.SetValue("M")

	return Model{
		viewMode:         ViewList,
		processes:        []core.Process{},
		selectedIdx:      0,
		statusFilter:     core.StatusRunning,
		textInput:        ti,
		spawnTitle:       titleInput,
		spawnDesc:        descInput,
		spawnPriority:    priorityInput,
		spawnFocusedField: 0,
	}
}

// Init initializes the TUI model
func (m Model) Init() tea.Cmd {
	// Load processes on init
	return tea.Batch(
		// Load initial processes with filter
		refreshProcessesWithFilter(m.statusFilter),
		// Start tick for auto-refresh
		tea.Tick(time.Second*5, func(t time.Time) tea.Msg { return TickMsg(t) }),
	)
}

// refreshProcesses returns a command to reload processes from DB
func refreshProcesses() tea.Cmd {
	return func() tea.Msg {
		processes, err := core.ListProcesses(nil)
		if err != nil {
			return errMsg{err}
		}
		return ProcessesLoadedMsg{Processes: processes}
	}
}

// refreshProcessesWithFilter returns a command to reload processes with a filter
func refreshProcessesWithFilter(filter core.ProcessStatus) tea.Cmd {
	return func() tea.Msg {
		var processes []core.Process
		var err error

		if filter == "" {
			processes, err = core.ListProcesses(nil)
		} else {
			processes, err = core.ListProcesses(&filter)
		}

		if err != nil {
			return errMsg{err}
		}
		return ProcessesLoadedMsg{Processes: processes}
	}
}

// errMsg wraps an error for Bubble Tea
type errMsg struct {
	error
}

// ProcessesLoadedMsg signals that processes have been loaded
type ProcessesLoadedMsg struct {
	Processes []core.Process
}
