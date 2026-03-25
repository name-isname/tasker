# taskctl

> A 3-in-1 task management tool with CLI, TUI, and Web interfaces

## Overview

**taskctl** is a task management tool written in Go that supports three different interfaces through a single binary distribution:

- **CLI** - Command-line interface for quick task operations
- **TUI** - Interactive terminal user interface with Bubble Tea
- **Web** - Browser-based UI with Vue.js

### Architecture

The project follows a **"Core-Library First"** design pattern:

```
┌─────────────────────────────────────────────────────┐
│                    core/                            │
│            (Business & Data Layer)                  │
│     AddTask() | ListTasks() | CompleteTask()        │
└──────────┬──────────────────────────────────────────┘
           │
     ┌─────┴─────┬──────────────┬─────────────┐
     ▼           ▼              ▼             ▼
  cli/        tui/            web/        (future)
(Cobra)    (Bubble Tea)      (Gin+Vue)
```

All database and business logic lives in `core/`. The presentation layers simply call core functions.

## Tech Stack

| Layer | Technology |
|-------|-----------|
| **Database** | SQLite via `github.com/glebarez/sqlite` (pure Go, no CGO) |
| **ORM** | GORM |
| **CLI** | Cobra |
| **TUI** | Bubble Tea + Lipgloss |
| **Web API** | Gin |
| **Frontend** | Vue 3 + TypeScript + Vite + TailwindCSS |

## Installation

### Prerequisites

- Go >= 1.22
- Node.js >= 18 (for frontend development)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/taskctl.git
cd taskctl

# Install frontend dependencies
make install-frontend

# Build everything (frontend + Go binary)
make build

# The binary will be created as ./taskctl
```

### Cross-Platform Build

```bash
make build-linux    # Linux AMD64
make build-mac      # macOS AMD64 + ARM64
make build-windows  # Windows AMD64
```

## Usage

### CLI

```bash
# Add a task
./taskctl add "Buy groceries"

# List all tasks
./taskctl list

# List tasks as JSON (for AI agents)
./taskctl list --json

# Mark task as completed
./taskctl complete 1

# Delete a task
./taskctl delete 1

# Use custom database location
./taskctl --db ~/mytasks.db list
```

### TUI (Terminal UI)

```bash
# Launch interactive terminal interface
./taskctl tui
```

Use `j/k` or `↑/↓` to navigate, `q` to quit.

### Web UI

```bash
# Start web server (default port: 8080)
./taskctl web

# Specify custom port
./taskctl web --port 3000
```

Then open http://localhost:8080 in your browser.

## Development

```bash
# Install frontend dependencies
make install-frontend

# Start frontend dev server
make dev

# Run CLI directly without building
make run

# Build only frontend
make build-frontend

# Build only Go binary
make build-go

# Clean build artifacts
make clean
```

## API Routes

The web server exposes RESTful API endpoints under `/api/v1/`:

- `GET /api/v1/tasks` - List all tasks
- `POST /api/v1/tasks` - Create a new task
- `PUT /api/v1/tasks/:id/complete` - Mark task as completed
- `DELETE /api/v1/tasks/:id` - Delete a task

## Project Structure

```
taskctl/
├── main.go              # Entry point
├── core/                # Business logic & data layer
│   ├── db.go           # Database initialization
│   └── task.go         # Task model and CRUD
├── cli/                 # CLI commands (Cobra)
│   ├── root.go         # Root command
│   ├── add.go          # Add task
│   ├── list.go         # List tasks (with --json)
│   ├── complete.go     # Complete task
│   ├── delete.go       # Delete task
│   ├── tui.go          # TUI entry
│   └── web.go          # Web server entry
├── tui/                 # Terminal UI (Bubble Tea)
│   ├── model.go        # TUI model
│   └── view.go         # TUI view rendering
├── web/                 # Web server & embedded frontend
│   ├── server.go       # Gin API server
│   ├── embed.go        # Static file embedding
│   └── frontend/       # Vue + Vite + TailwindCSS
│       ├── src/
│       ├── dist/       # Built assets (embedded)
│       └── ...
├── Makefile            # Build commands
├── CLAUDE.md           # AI assistant guide
└── README.md           # This file
```

## License

MIT License
