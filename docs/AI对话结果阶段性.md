


这是一份为你深度压缩的**“项目核心规格与设计文档”**。它剔除了我们讨论中的客套话，提取了所有的**架构决策、业务哲学和数据模型**。

你可以直接复制以下内容，带入到新的对话中（无论是发给 Claude Code 还是其他的 AI 助手），AI 瞬间就能完整理解你的意图并开始编码：

***

# Project Spec: Process-Oriented Task Manager (CLI + TUI + Web)

## 1. Core Philosophy
This is NOT a traditional GTD/Deadline-driven To-Do app. It is a **"Process Management"** tool designed to reduce human working memory burden. 
- Tasks are treated as OS **Processes**. 
- Focus is on **state transitions** (running, blocked, suspended) rather than just "completed".
- **Process = Result**: Every incremental progress or roadblock is recorded as a chronical **Log**, turning the task tracker into a personal knowledge base (PKM).

## 2. Architecture: "Core-Library First" & Single Binary
- **Language**: Go (>= 1.22).
- **Format**: Single Binary containing everything.
- **Data Layer**: `gorm` + `glebarez/sqlite` (Pure Go SQLite, strictly NO CGO for easy cross-compilation).
- **3-in-1 Interfaces**:
  1. **Core Layer**: Independent business logic and DB operations (No UI logic allowed here).
  2. **CLI (`spf13/cobra`)**: AI-First design. Must support `--json` flag and strict exit codes for LLM Agents (e.g., Claude Code).
  3. **TUI (`charmbracelet/bubbletea`)**: Keyboard-driven interactive terminal UI for geeks.
  4. **Web (`gin` + Vue3/Vite)**: Embedded via `//go:embed`. Serves RESTful API (`/api/v1/`) and static frontend.

## 3. Database Schema (Entities)

### Entity 1: `Process` (The Task/Goal)
| Field | Type | Note |
| :--- | :--- | :--- |
| `id` | INTEGER | Primary Key (Auto Increment for easy CLI usage) |
| `parent_id` | INTEGER | Foreign key to self (Supports infinite nested sub-processes) |
| `title` | VARCHAR | Short name of the process |
| `description`| TEXT | Detailed context (Markdown supported) |
| `status` | VARCHAR | `running`, `blocked`, `suspended`, `terminated` |
| `priority` | VARCHAR | `low`, `medium`, `high` |
| `ranking` | REAL | Custom sort weight (like Linux `nice` value) |
| `created_at` / `updated_at` | DATETIME | Auto-managed |

### Entity 2: `Log` (The Timeline/Memory)
| Field | Type | Note |
| :--- | :--- | :--- |
| `id` | INTEGER | Primary Key |
| `process_id` | INTEGER | Indexed, Links to Process |
| `log_type` | VARCHAR | `state_change` (auto-generated) or `progress` (user notes) |
| `content` | TEXT | Markdown text recording ideas, roadblocks, or progress |
| `created_at` | DATETIME | Timestamp of the log |

## 4. Key Technical Features Required
1. **Infinite Sub-processes**: Supported via `parent_id` (Queryable via SQLite `WITH RECURSIVE` if needed).
2. **Full-Text Search (FTS5)**: Must use SQLite FTS5 extension to enable millisecond global search across `Process.title`, `Process.description`, and `Log.content`.
3. **Markdown Export**: Provide a feature/CLI command to export a Process and its entire Log timeline as a single well-formatted `.md` file for PKM storage.

*** 

### 💡 附带的“新对话启动指令”建议：
把上面的英文发给 AI 后，你可以紧接着补充一句：
> *"Please act as a senior Go/Vue architect. Read the specification above. Start by initializing the Go module and implementing the `core/` package containing the GORM models and SQLite FTS5 initialization. Let me know when the core database layer is done."*
