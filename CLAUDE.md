# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture

**Core-Library First**: The `core/` package contains ALL business logic and data access. The `cli/`, `tui/`, and `web/` packages are pure presentation layers that MUST call functions from `core/`. Never write raw SQL, GORM queries, or business logic in presentation layers.

```
core/    → CreateProcess(), ListProcesses(), ChangeProcessState(), DeleteProcess(), AddLog(), GlobalSearch()
cli/     → Cobra commands that call core functions
tui/     → Bubble Tea UI that calls core functions
web/     → Gin API handlers that call core functions
```

**Single Binary Distribution**: The Vue frontend is built to `web/frontend/dist` and embedded into the Go binary via `//go:embed` in `web/embed.go`. This allows distribution as a single executable.

**Database Initialization**: The Cobra root command has a `PersistentPreRunE` hook that automatically initializes the SQLite database and runs migrations before ANY subcommand executes. The DB path is configurable via `--db` flag (default: `./taskctl.db`).

## Data Model

**Process-Oriented Task Management**: Tasks are modeled as OS "Processes" with state transitions rather than simple todo items. This is a key conceptual distinction - tasks have lifecycle states (running, blocked, suspended, terminated) and can have child sub-processes.

### Process Entity
- `id`: Primary key (auto-increment)
- `parent_id`: Self-referencing foreign key for infinite nested sub-processes
- `title`: Short name
- `description`: Detailed context (Markdown)
- `status`: `running`, `blocked`, `suspended`, `terminated`
- `priority`: `low`, `medium`, `high`
- `ranking`: Custom sort weight (like Linux `nice` value)

### Log Entity
- `id`: Primary key
- `process_id`: Foreign key to Process
- `log_type`: `state_change` (auto-created by `ChangeProcessState`) or `progress` (manual)
- `content`: Markdown text recording ideas, roadblocks, or progress

### FTS5 Full-Text Search
SQLite FTS5 virtual table `process_fts` with auto-sync triggers for millisecond global search across processes and logs.

## Build Commands

```bash
make build           # Build frontend (npm) then Go binary
make build-frontend  # Only build Vue frontend to web/frontend/dist
make build-go        # Only compile Go binary (taskctl)
make dev             # Start Vite dev server on :5173
make clean           # Remove build artifacts
make install-frontend # Install npm dependencies
./taskctl            # Run the CLI directly
go test ./...        # Run all tests
go test ./tui/...    # Run TUI tests only
```

Cross-platform builds: `make build-linux`, `make build-mac`, `make build-windows`

## Technology Constraints

- **SQLite**: MUST use `github.com/glebarez/sqlite` (pure Go, no CGO) to enable easy cross-compilation
- **CLI**: `taskctl ps` MUST support `--json` flag for AI agent compatibility
- **Web API**: All API routes are prefixed with `/api/v1/`
- **Database location**: Configurable via `--db` flag (default: `./taskctl.db`)

## Frontend

Vue 3 + TypeScript + Vite + TailwindCSS. The frontend build output must be in `web/frontend/dist` for Go embedding.

## TUI Implementation Notes

The TUI uses Bubble Tea with a ViewMode enum pattern to manage different screens. Key implementation details:

- **ViewMode enum**: ViewList, ViewDetail, ViewInput, ViewHelp, ViewSpawn, ViewEditProcess, ViewSearch, ViewTimeline, ViewStats, ViewTree, ViewParentSelect
- **Message handling**: Separate handler functions for each ViewMode (e.g., `handleListKeyMsg`, `handleDetailKeyMsg`)
- **Multi-line input**: Use `bubbles/textarea` for description and log fields, `bubbles/textinput` for single-line fields
- **Ctrl+Enter for submission**: Forms with textarea fields use Ctrl+Enter to submit, regular Enter allows line breaks
- **Parent process selection**: Separate ViewParentSelect mode for choosing parent processes
- **Viewport scrolling**: List views use viewportOffset for large datasets

## Core Functions Reference

Process Operations (in `core/process.go`):
- `CreateProcess(title, description, parentID, priority)` - Create new process
- `GetProcess(id)` - Get single process
- `ListProcesses(filter)` - List processes with optional status filter
- `UpdateProcess(id, title, desc, priority)` - Update process fields
- `DeleteProcess(id)` - Delete process and all descendants
- `ChangeProcessState(id, newStatus, reason)` - Atomically change status and create state-change log

Log Operations (in `core/log.go`):
- `AddLog(processID, logType, content)` - Add a log entry
- `GetLogs(processID)` - Get all logs for a process
- `UpdateLog(id, content)` - Update log content
- `DeleteLog(id)` - Delete a log

Advanced Features (in `core/advanced.go`):
- `GlobalSearch(keyword)` - FTS5 search across processes and logs
- `GetTimeline(startTime, endTime, limit)` - Get chronological activity
- `GetTodayTimeline()` - Get today's activity
- `GetActivityStats(days)` - Get activity counts per day
- `GetProcessTree(rootID)` - Get process with descendants
- `GetFullProcessTree()` - Get all root processes with trees
- `ExportProcessMarkdown(id)` - Export process + logs as Markdown

## CLI Commands

Process Management:
- `taskctl spawn <title>` - Create new process
- `taskctl ps` - List processes (--json for AI)
- `taskctl inspect <id>` - Show process details
- `taskctl update <id>` - Edit process
- `taskctl kill <id>` - Delete process
- `taskctl block <id>` - Set to blocked state
- `taskctl wake <id>` - Set to running state
- `taskctl terminate <id>` - Set to terminated state

Log Management:
- `taskctl log add <process-id> <content>` - Add log
- `taskctl logs <process-id>` - List process logs
- `taskctl log update <log-id> <content>` - Update log
- `taskctl log rm <log-id>` - Delete log

Advanced:
- `taskctl grep <keyword>` - Global search
- `taskctl timeline` - Show activity timeline
- `taskctl stats [days]` - Show activity statistics
- `taskctl tree` - Show process tree
- `taskctl export <id>` - Export as Markdown

Interface:
- `taskctl tui` - Launch terminal UI
- `taskctl web` - Start web server
