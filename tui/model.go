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
	ViewSearch
	ViewTimeline
	ViewStats
	ViewTree
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

	// Search view state
	searchKeyword  string
	searchResults  []core.SearchResult
	searchCursor   int

	// Timeline view state
	timelineEntries []core.TimelineEntry
	timelineCursor  int

	// Stats view state
	statsDays      int
	statsData      []core.ActivityStat

	// Tree view state
	treeNodes      []*core.ProcessNode
	treeCursor     int
	treeExpanded   map[uint]bool // Track expanded nodes
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
		viewMode:          ViewList,
		processes:         []core.Process{},
		selectedIdx:       0,
		statusFilter:      core.StatusRunning,
		textInput:         ti,
		spawnTitle:        titleInput,
		spawnDesc:         descInput,
		spawnPriority:     priorityInput,
		spawnFocusedField: 0,
		treeExpanded:      make(map[uint]bool),
		statsDays:         30,
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

// SearchResultsLoadedMsg signals that search results have been loaded
type SearchResultsLoadedMsg struct {
	Keyword  string
	Results  []core.SearchResult
}

// TimelineLoadedMsg signals that timeline has been loaded
type TimelineLoadedMsg struct {
	Entries []core.TimelineEntry
}

// StatsLoadedMsg signals that stats have been loaded
type StatsLoadedMsg struct {
	Days int
	Stats []core.ActivityStat
}

// TreeLoadedMsg signals that tree has been loaded
type TreeLoadedMsg struct {
	Nodes []*core.ProcessNode
}

// loadSearchResults returns a command to search for a keyword
func loadSearchResults(keyword string) tea.Cmd {
	return func() tea.Msg {
		results, err := core.GlobalSearch(keyword)
		if err != nil {
			return errMsg{err}
		}
		return SearchResultsLoadedMsg{Keyword: keyword, Results: results}
	}
}

// loadTimeline returns a command to load the timeline
func loadTimeline(days int) tea.Cmd {
	return func() tea.Msg {
		var entries []core.TimelineEntry
		var err error

		if days > 0 {
			startTime := time.Now().AddDate(0, 0, -days)
			entries, err = core.GetTimeline(startTime, time.Time{}, 100)
		} else {
			entries, err = core.GetTodayTimeline()
		}

		if err != nil {
			return errMsg{err}
		}
		return TimelineLoadedMsg{Entries: entries}
	}
}

// loadStats returns a command to load activity stats
func loadStats(days int) tea.Cmd {
	return func() tea.Msg {
		stats, err := core.GetActivityStats(days)
		if err != nil {
			return errMsg{err}
		}
		return StatsLoadedMsg{Days: days, Stats: stats}
	}
}

// loadTree returns a command to load the process tree
func loadTree() tea.Cmd {
	return func() tea.Msg {
		nodes, err := core.GetFullProcessTree()
		if err != nil {
			return errMsg{err}
		}
		return TreeLoadedMsg{Nodes: nodes}
	}
}
