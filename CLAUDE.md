# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture

**Core-Library First**: The `core/` package contains ALL business logic and data access. The `cli/`, `tui/`, and `web/` packages are pure presentation layers that MUST call functions from `core/`. Never write raw SQL, GORM queries, or business logic in presentation layers.

```
core/    → CreateProcess(), ListProcesses(), SetProcessStatus(), DeleteProcess(), AddLog(), SearchProcesses()
cli/     → Cobra commands that call core functions
tui/     → Bubble Tea UI that calls core functions
web/     → Gin API handlers that call core functions
```

**Single Binary Distribution**: The Vue frontend is built to `web/frontend/dist` and embedded into the Go binary via `//go:embed` in `web/embed.go`. This allows distribution as a single executable.

**Database Initialization**: The Cobra root command has a `PersistentPreRunE` hook that automatically initializes the SQLite database and runs migrations before ANY subcommand executes. The DB path is configurable via `--db` flag (default: `./taskctl.db`).

## Data Model

**Process-Oriented Task Management**: Tasks are modeled as OS "Processes" with state transitions rather than simple todo items.

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
- `log_type`: `state_change` (auto) or `progress` (manual)
- `content`: Markdown text recording ideas, roadblocks, or progress

### FTS5 Full-Text Search
SQLite FTS5 virtual table `process_fts` with auto-sync triggers for millisecond global search.

## Build Commands

```bash
make build           # Build frontend (npm) then Go binary
make build-frontend  # Only build Vue frontend to web/frontend/dist
make build-go        # Only compile Go binary (taskctl)
make dev             # Start Vite dev server on :5173
make clean           # Remove build artifacts
make install-frontend # Install npm dependencies
./taskctl            # Run the CLI directly
```

Cross-platform builds: `make build-linux`, `make build-mac`, `make build-windows`

## Technology Constraints

- **SQLite**: MUST use `github.com/glebarez/sqlite` (pure Go, no CGO) to enable easy cross-compilation
- **CLI**: `taskctl list` MUST support `--json` flag for AI agent compatibility
- **Web API**: All API routes are prefixed with `/api/v1/`
- **Database location**: Configurable via `--db` flag (default: `./taskctl.db`)

## Frontend

Vue 3 + TypeScript + Vite + TailwindCSS. The frontend build output must be in `web/frontend/dist` for Go embedding.

## Implementation Status

**Completed**: Core data models (Process, Log), CLI commands (add, list, status, delete), FTS5 full-text search, basic project scaffolding

**Incomplete / TODO**:
- `tui/model.go` and `tui/view.go` - Bubble Tea model and view are skeleton only
- `cli/tui.go` - TUI entry is a stub, needs `tea.NewProgram()` implementation
- Markdown export functionality - Export Process + Logs as `.md` file
- Frontend UI components - still using Vite default template, needs task management UI

## Key Files

- `main.go` → Entry point, calls `cli.Execute()`
- `core/models.go` → Process and Log GORM models
- `core/db.go` → GORM + SQLite initialization with FTS5 setup
- `core/process.go` → Process CRUD operations
- `core/log.go` → Log CRUD operations
- `cli/root.go` → Cobra root with PersistentPreRunE for DB init
- `web/server.go` → Gin server with `/api/v1/` routes
- `web/embed.go` → `//go:embed frontend/dist` for static file serving
