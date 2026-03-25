# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture

**Core-Library First**: The `core/` package contains ALL business logic and data access. The `cli/`, `tui/`, and `web/` packages are pure presentation layers that MUST call functions from `core/`. Never write raw SQL, GORM queries, or business logic in presentation layers.

```
core/    → AddTask(), ListTasks(), CompleteTask(), DeleteTask()
cli/     → Cobra commands that call core functions
tui/     → Bubble Tea UI that calls core functions
web/     → Gin API handlers that call core functions
```

**Single Binary Distribution**: The Vue frontend is built to `web/frontend/dist` and embedded into the Go binary via `//go:embed` in `web/embed.go`. This allows distribution as a single executable.

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

## Key Files

- `main.go` → Entry point, calls `cli.Execute()`
- `core/db.go` → GORM + SQLite initialization
- `core/task.go` → Task model and CRUD exports
- `cli/root.go` → Cobra root with PersistentPreRunE for DB init
- `web/server.go` → Gin server with `/api/v1/` routes
- `web/embed.go` → `//go:embed frontend/dist` for static file serving
