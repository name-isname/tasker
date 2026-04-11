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

**Git Workflow**: This project has automated hooks (`.claude/settings.json`):
- **PostToolUse (Edit|Write)**: Auto-analyzes changes and creates conventional commits via agent
- **Stop**: Runs `make build` when Claude Code session ends

ALWAYS create a git commit after completing code changes. Use conventional commit format:
- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation changes
- `refactor:` for code refactoring
- `test:` for test changes
- `chore:` for maintenance tasks

## Core Package Structure

The `core/` package is organized into:
- `models.go` - Entity definitions (Process, Log, ProcessFTS) and constants
- `db.go` - Database initialization and FTS5 full-text search setup
- `process.go` - Process CRUD operations
- `log.go` - Log CRUD and pagination
- `advanced.go` - Tree operations, timeline, stats, search, export
- `dto.go` - Data transfer objects (ProcessNode, TimelineEntry, ActivityStat, SearchResult, ProcessExport)
- `testutil.go` - Test database setup utilities

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

### Circular Reference Prevention
The system prevents circular parent-child relationships at write time and handles existing corrupted data at read time:
- `CreateProcess` validates that the specified parent exists
- `UpdateProcess` uses `wouldCreateCircularReference()` to detect cycles before applying changes
- Recursive operations (`DeleteProcess`, `GetProcessTree`, `GetProcessContext`, `GetFullProcessTree`) include cycle detection as a safety measure
- Returns `ErrCircularReference` when a cycle would be created

## Build Commands

```bash
make build           # Build frontend (npm) then Go binary
make build-frontend  # Only build Vue frontend to web/frontend/dist
make build-go        # Only compile Go binary (taskctl)
make dev             # Start Vite dev server on :5173
make clean           # Remove build artifacts
make install-frontend # Install npm dependencies
make run             # Run CLI directly without building
make release-snapshot # Test goreleaser release (does not publish)
make release         # Execute goreleaser release (requires GITHUB_TOKEN)
./taskctl            # Run the built binary
go test ./...        # Run all tests
go test ./core/...   # Run core package tests only
go test -run TestFoo # Run specific test
```

## Release Process

This project uses **goreleaser** for automated cross-platform releases.

### Pre-release Checklist

1. **Ensure all changes are committed** - goreleaser requires a clean git state
2. **Run tests** - Verify all tests pass: `go test ./...`
3. **Update version** - Decide on new version number (semantic versioning)

### Release Steps

```bash
# 1. Test release build (does not publish, validates configuration)
goreleaser release --snapshot --clean

# 2. Clean snapshot artifacts
rm -rf dist/

# 3. Create and push git tag
git tag -a v0.2.0 -m "Release v0.2.0: description of changes"
git push origin main        # Push commits first
git push origin v0.2.0      # Then push the tag

# 4. Execute official release (requires GITHUB_TOKEN)
export GITHUB_TOKEN="ghp_xxx"  # GitHub PAT with repo permissions
goreleaser release

# 5. Manually update Homebrew cask (goreleaser v2 has directory issues)
./scripts/update-homebrew.sh v0.2.0
```

**Important**:
- Push commits BEFORE pushing the tag
- goreleaser will fail if git state is dirty (uncommitted changes)
- GITHUB_TOKEN must have `repo` scope permissions
- Homebrew cask MUST be manually updated after release (see below)

### What Goreleaser Does

1. **Builds binaries** for:
   - macOS (amd64, arm64)
   - Linux (amd64, arm64)
   - Windows (amd64)

2. **Creates GitHub Release** with:
   - All build artifacts as downloadable files
   - checksums.txt with SHA256 hashes
   - Auto-generated changelog

3. **Generates Homebrew cask** locally in `dist/homebrew/taskctl.rb`
   - Note: goreleaser v2 has issues pushing to correct directory
   - Use the provided script to manually update the cask

### Version Information

Version info is injected via ldflags during goreleaser builds:
- `main.version` - Version string from git tag (e.g., "v0.2.0")
- `main.commit` - Git commit SHA
- `main.date` - Build timestamp

Access via: `taskctl version`

### Homebrew Distribution

The project is distributed via Homebrew tap at `name-isname/homebrew-taskctl`.

**Installation**:
```bash
brew install --cask name-isname/taskctl/taskctl
```

**Upgrade after release**:
```bash
brew upgrade --cask name-isname/taskctl/taskctl
```

**Cask location**: `/opt/homebrew/Library/Taps/name-isname/homebrew-taskctl/Casks/taskctl.rb`

**Manual Cask Update** (after goreleaser release):
```bash
# Use the provided script to update homebrew cask
./scripts/update-homebrew.sh v0.2.0

# Or manually update:
# 1. Get checksums from the release
curl -sL https://github.com/name-isname/tasker/releases/download/v0.2.0/tasker_0.2.0_darwin_arm64.tar.gz | shasum -a 256
# (repeat for each platform)

# 2. Clone and update the tap
cd /tmp
git clone https://github.com/name-isname/homebrew-taskctl.git
cd homebrew-taskctl

# 3. Edit Casks/taskctl.rb with new version, URLs, and checksums

# 4. Commit and push
git add Casks/taskctl.rb
git commit -m "Update taskctl cask to v0.2.0"
git push
```

**Why manual update?** goreleaser v2 has a known issue where it pushes cask to the repository root (`/taskctl.rb`) instead of the correct `Casks/` directory. The provided script handles this correctly.

**Cask supports**: macOS ARM64/AMD64, Linux ARM64/AMD64

**macOS Quarantine Removal**: The goreleaser config includes a post-install hook (`homebrew_casks.hooks.post.install`) that automatically removes the macOS quarantine attribute from the downloaded binary using `/usr/bin/xattr -dr com.apple.quarantine`. This prevents macOS from showing the "unidentified developer" warning on first run.

**Changelog Generation**: goreleaser auto-generates release notes from commits, excluding `docs:`, `test:`, and `chore:` prefixes. Only `feat:`, `fix:`, and `refactor:` commits appear in the GitHub release notes.

## Technology Constraints

- **SQLite**: MUST use `github.com/glebarez/sqlite` (pure Go, no CGO) to enable easy cross-compilation. Goreleaser builds set `CGO_ENABLED=0`.
- **CLI**: `taskctl ps` MUST support `--xml` flag for AI agent compatibility
- **Web API**: All API routes are prefixed with `/api/v1/`
- **Database location**: Configurable via `--db` flag (default: `./taskctl.db`)
- **Status changes**: MUST use `ChangeProcessState()` (transactional) instead of `SetProcessStatus()` for consistency
- **Version info**: Build-time variables (`main.version`, `main.commit`, `main.date`) are injected via goreleaser ldflags

## Frontend Architecture

**Stack**: Vue 3.5 (Composition API), TypeScript 5.9, Vite 8, TailwindCSS 4.2, Vue Router, Pinia, Axios

**Structure**: The frontend follows a standard Vue 3 composition API pattern with centralized state management:

```
web/frontend/src/
├── types/api.ts          # TypeScript interfaces for API contracts
├── services/api.ts       # Axios client with /api/v1 base URL
├── stores/
│   ├── processes.ts      # Process list state (Pinia)
│   └── theme.ts          # Dark mode state with localStorage persistence
├── composables/
│   ├── useTheme.ts       # Theme toggle composable (light/dark/system)
│   └── useKeyboard.ts    # Global keyboard shortcuts handler
├── router/index.ts       # Vue Router (/, /process/:id, /search)
├── components/           # Reusable UI components
└── views/                # Page-level components
```

**Frontend Development Workflow**:
```bash
cd web/frontend
npm install               # Install dependencies
npm run dev              # Start Vite dev server on :5173 (proxying /api to backend)
npm run build            # Build to dist/ for Go embedding
```

**Dark Mode**: Three-state system (light/dark/system) with class-based toggle. The `dark` class on `<html>` controls Tailwind's dark mode variant.

**Keyboard Shortcuts**:
- Global: `c` (create), `/` (search), `n` (home), `?` (help), `Escape` (close)
- List view: `j/k` or `↑/↓` (navigate), `Enter` (open), `d` (delete), `e` (edit)
- Detail view: `E` (edit process), `b/p/w/t` (status change), `a` (add log), `e` (edit log), `x` (delete log), `X` (export Markdown), `m` (toggle Markdown), `j/k` (select log), `1-4` (quick status change)

**Production**: Built assets are embedded via `//go:embed` in `web/embed.go`. The Gin server uses `NoRoute` handler to serve the SPA, returning `index.html` for all non-API routes.

**Important**: When modifying the frontend, always test both light and dark modes. Components use CSS variables (`--bg`, `--text`, `--border`, `--accent`) for theming.

## TUI Implementation Notes

The TUI uses Bubble Tea with a ViewMode enum pattern to manage different screens. Key implementation details:

- **Markdown rendering**: `tui/markdown.go` provides terminal-friendly markdown rendering for descriptions and logs (supports bold, italic, code, headers, lists, quotes, code blocks)

- **ViewMode enum**: ViewList, ViewDetail, ViewInput, ViewHelp, ViewSpawn, ViewEditProcess, ViewSearch, ViewTimeline, ViewStats, ViewTree, ViewParentSelect, ViewDeleteConfirm, ViewExportConfirm
- **Message handling**: Separate handler functions for each ViewMode (e.g., `handleListKeyMsg`, `handleDetailKeyMsg`)
- **Multi-line input**: Use `bubbles/textarea` for description and log fields, `bubbles/textinput` for single-line fields
- **Ctrl+Enter for submission**: Forms with textarea fields use Ctrl+Enter to submit, regular Enter allows line breaks
- **Parent process selection**: Separate ViewParentSelect mode for choosing parent processes; filters out current process AND all descendants
- **Viewport scrolling**: List views use viewportOffset for large datasets
- **Auto-refresh behavior**: TickMsg only triggers refresh when in ViewList mode to prevent kicking users out of detail/search/timeline views

### TUI Message Types
The TUI uses custom Bubble Tea messages for async operations:
- `ProcessesLoadedMsg` - Standard process list refresh (cursor reset if out of bounds)
- `ProcessDeletedMsg` - Post-deletion refresh with `DeletedIndex` for smart cursor adjustment
- `ProcessDetailLoadedMsg` - Process detail view with logs
- `ShowDetailMsg`, `BackToListMsg` - Navigation messages
- `ParentsLoadedMsg` - Parent selection list loaded
- `ExportSuccessMsg` - Export completed with file path
- `ClearExportSuccessMsg` - Auto-clear export success message after delay
- `errMsg` - Error handling

### TUI Export Functionality
Export is implemented in `tui/export.go`:
- `ExportProcess(processID)` - Exports process to Markdown file in current directory
- `GenerateExportFileName(process)` - Generates `{title}-{id}.md` filename
- `GetAbsolutePath(filename)` - Resolves relative paths to absolute
- Detail view shortcut: `X` (uppercase) triggers export confirmation dialog
- Success message auto-clears after 3 seconds using `tea.Tick`

## Feature Implementation Guidelines

When implementing deletion features, ensure:
1. database record is removed
2. list view refreshes to show updated state
3. confirmation uses 'y' key
4. cursor is adjusted to point to the next item (or last item if deleting the last one) - use `ProcessDeletedMsg` instead of `ProcessesLoadedMsg` for proper cursor handling

When implementing parent selection:
1. Filter out the current process being edited
2. Filter out ALL descendants of the current process (use `GetDescendantIDs()`)
3. Display cached parent name for better UX

## Core Functions Reference

Process Operations (in `core/process.go`):
- `CreateProcess(title, description, parentID, priority)` - Create new process (validates parent exists)
- `GetProcess(id)` - Get single process
- `ListProcesses(filter)` - List processes with optional status filter
- `UpdateProcess(id, title, desc, priority, parentID)` - Update process fields (detects circular references)
- `DeleteProcess(id)` - Delete process and all descendants (with cycle detection)
- `ChangeProcessState(id, newStatus, reason)` - Atomically change status and create state-change log (preferred over SetProcessStatus)
- `SetProcessStatus(id, status)` - Legacy method, use ChangeProcessState instead
- `SetProcessRanking(id, ranking)` - Update sort weight
- `GetChildProcesses(parentID)` - Get direct children
- `GetRootProcesses()` - Get top-level processes
- `GetDescendantIDs(parentID)` - Get all descendant IDs (useful for filtering)
- `wouldCreateCircularReference(processID, parentID)` - Check if parent relationship would create a cycle

Log Operations (in `core/log.go`):
- `AddLog(processID, logType, content)` - Add a log entry
- `GetLogs(processID)` - Get all logs for a process
- `GetLogsPaginated(processID, page, pageSize)` - Get paginated logs with total count
- `UpdateLog(id, content)` - Update log content
- `DeleteLog(id)` - Delete a log
- `GetAllLogs(logType, limit)` - Get logs across all processes

Advanced Features (in `core/advanced.go`):
- `ChangeProcessState(id, newStatus, reason)` - Transactional status change with log
- `GetProcessTree(rootID)` - Get process with descendants (cycle-safe)
- `GetFullProcessTree()` - Get all root processes with trees (cycle-safe)
- `GetProcessContext(processID)` - Get process with logs and children for export (cycle-safe)
- `GetTimeline(startTime, endTime, limit)` - Get chronological activity
- `GetTodayTimeline()` - Get today's activity
- `GetActivityStats(days)` - Get activity counts per day
- `GlobalSearch(keyword)` - Search across processes and logs
- `GetActiveProcesses()` - Get running processes with recent logs
- `GetBlockedProcesses()` - Get all blocked processes
- `FormatProcessTree(node, prefix, isLast)` - Render tree as ASCII
- `FormatFullTree()` - Render entire forest as ASCII

**DTO Types** (in `core/dto.go`):
- `ProcessNode` - Hierarchical tree structure for visualization
- `TimelineEntry` - Log entry with parent process context
- `ActivityStat` - Daily log counts for heatmaps
- `SearchResult` - Unified search result with type indicator
- `ProcessExport` - Complete process data for Markdown export

## CLI Commands

Utility:
- `taskctl version` - Print version, commit, and build date information

Process Management:
- `taskctl spawn <title>` - Create new process
- `taskctl ps` - List processes (--xml for AI)
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

## Web API

All routes prefixed with `/api/v1/`:
- `GET /api/v1/processes` - List processes (optional `?status=` filter)
- `GET /api/v1/processes/:id` - Get single process
- `POST /api/v1/processes` - Create process
- `PUT /api/v1/processes/:id/status` - Change status (uses `ChangeProcessState` transactionally)
- `DELETE /api/v1/processes/:id` - Delete process
- `GET /api/v1/processes/:id/logs` - Get process logs
- `POST /api/v1/processes/:id/logs` - Add log
- `GET /api/v1/search?q=` - Global search

## Testing

Tests are organized by package:
- `core/process_test.go` - Process CRUD and circular reference tests
- `core/log_test.go` - Log CRUD and pagination tests
- `core/advanced_test.go` - Advanced features including cycle-safe operations
- `tui/model_test.go` - TUI model tests
- `cli/cli_test.go` - CLI command tests

Run tests with:
```bash
go test ./...              # All tests
go test ./core/...         # Core package only
go test -run TestFoo       # Specific test
go test -v ./core/...      # Verbose output
```

Test utilities in `core/testutil.go` provide `setupTestDB()` and `teardownTestDB()` for isolated test databases.
