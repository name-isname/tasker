package tui

import (
	"taskctl/core"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
	ViewParentSelect  // For selecting parent process
	ViewDeleteConfirm // For delete confirmation
	ViewExportConfirm // For export confirmation
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

	// Help page scroll offset
	helpOffset int

	// Status filter
	statusFilter core.ProcessStatus
	filtering    bool

	// Text input
	textInput          textarea.Model
	inputPrompt        string
	editingLogID       uint                // 0 for new log, >0 for editing existing log
	pendingStateChange *core.ProcessStatus // Pending state change with note

	// Spawn/Edit form fields (shared between spawn and edit)
	spawnTitle        textinput.Model
	spawnDesc         textarea.Model
	spawnPriority     textinput.Model
	spawnFocusedField int  // 0=title, 1=desc, 2=priority, 3=parent
	editingProcessID  uint // 0 for new process, >0 for editing existing process

	// Parent process selection
	availableParents   []core.Process // List of processes that can be parents
	parentCursor       int            // Cursor for parent selection
	selectedParentID   *uint          // Currently selected parent ID (nil = no parent)
	selectedParentName string         // Cached parent title for display

	// Log selection in detail view
	logCursor int // Index of selected log in processLogs

	// Search view state
	searchKeyword string
	searchResults []core.SearchResult
	searchCursor  int

	// Timeline view state
	timelineEntries []core.TimelineEntry
	timelineCursor  int

	// Stats view state
	statsDays int
	statsData []core.ActivityStat

	// Tree view state
	treeNodes    []*core.ProcessNode
	treeCursor   int
	treeExpanded map[uint]bool // Track expanded nodes

	// Delete confirmation state
	confirmDeleteType  string // "process" or "log"
	confirmDeleteID    uint   // ID of item to delete
	confirmDeleteName  string // Name/title of item to delete
	confirmDeleteIndex int    // Index of item in list before deletion (for cursor adjustment)

	// Export state
	exportProcessID  uint   // ID of process to export
	exportFileName   string // Generated filename for export
	exportFilePath   string // Full path where file will be written
	exportSuccessMsg string // Success message after export

	// Debug: show last pressed key
	lastKey string

	// Markdown rendering state
	markdownEnabled     bool // Toggle markdown rendering on/off
	viewingRawMarkdown  bool // Temporarily view raw markdown text
}

// Messages
type (
	TickMsg       time.Time
	RefreshMsg    struct{}
	ShowDetailMsg struct {
		ProcessID uint
	}
	BackToListMsg   struct{}
	StatusChangeMsg struct {
		ProcessID uint
		Status    core.ProcessStatus
	}
	ProcessDeletedMsg struct {
		Processes      []core.Process
		DeletedIndex   int // Index of deleted item before deletion
	}
	ExportSuccessMsg struct {
		FilePath string
	}
)

// InitialModel creates the initial TUI model
func InitialModel() Model {
	ti := textarea.New()
	ti.Placeholder = "输入日志内容（支持多行，Cmd+Enter确认）..."
	ti.Focus()

	titleInput := textinput.New()
	titleInput.Placeholder = "进程标题..."

	descInput := textarea.New()
	descInput.Placeholder = "描述（可选，支持多行）..."

	priorityInput := textinput.New()
	priorityInput.Placeholder = "M"
	priorityInput.SetValue("M")

	return Model{
		viewMode:          ViewList,
		processes:         []core.Process{},
		selectedIdx:       0,
		statusFilter:      "", // Show all by default
		textInput:         ti,
		spawnTitle:        titleInput,
		spawnDesc:         descInput,
		spawnPriority:     priorityInput,
		spawnFocusedField: 0,
		treeExpanded:      make(map[uint]bool),
		statsDays:         30,
		availableParents:  []core.Process{},
		parentCursor:      -1,
		selectedParentID:  nil,
		markdownEnabled:   true, // Enable markdown by default
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
	Keyword string
	Results []core.SearchResult
}

// TimelineLoadedMsg signals that timeline has been loaded
type TimelineLoadedMsg struct {
	Entries []core.TimelineEntry
}

// StatsLoadedMsg signals that stats have been loaded
type StatsLoadedMsg struct {
	Days  int
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

// ParentsLoadedMsg signals that available parents have been loaded
type ParentsLoadedMsg struct {
	Parents []core.Process
}

// loadAvailableParents returns a command to load available parent processes
func loadAvailableParents() tea.Cmd {
	return func() tea.Msg {
		parents, err := core.ListProcesses(nil)
		if err != nil {
			return errMsg{err}
		}
		return ParentsLoadedMsg{Parents: parents}
	}
}
