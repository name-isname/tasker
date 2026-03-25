package tui

import (
	"time"
	"github.com/charmbracelet/bubbletea"
	"taskctl/core"
)

// ViewMode represents the current view state
type ViewMode int

const (
	ViewList ViewMode = iota
	ViewDetail
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
	return Model{
		viewMode:     ViewList,
		processes:    []core.Process{},
		selectedIdx:  0,
		statusFilter: core.StatusRunning,
	}
}

// Init initializes the TUI model
func (m Model) Init() tea.Cmd {
	// Load processes on init
	return tea.Batch(
		// Load initial processes
		refreshProcesses(),
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

// errMsg wraps an error for Bubble Tea
type errMsg struct {
	error
}

// ProcessesLoadedMsg signals that processes have been loaded
type ProcessesLoadedMsg struct {
	Processes []core.Process
}
