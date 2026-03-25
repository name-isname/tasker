


你提醒得非常对！我们之前的设计确实更偏向于“状态流转（高级功能）”，而忽略了最基本的**直接修改（Update）**和**单体查看（Read-One）**，以及针对日志本身改错的 CRUD。

在一个本地工具中，拼写错误是常有的事，所以我们必须提供直接修改标题、描述以及修改单条日志的能力。

为了保持我们“操作系统”的极客隐喻，我将这些基础 CRUD 命令融入了体系，比如用 `inspect` 来查看单个进程详情，用 `update` 来修改属性，并补全了针对 Log 的 CRUD。

以下是补充了基础 CRUD 之后的**完整修订版 `CLI_SPEC.md`**：

***

# Process Management - CLI Specification (Complete Edition)

## 1. Overview & CLI Philosophy
This document defines the Command Line Interface (CLI) for `taskctl`. 
The CLI follows a Unix-like process management metaphor. It is built using `github.com/spf13/cobra`.

**Crucial Rules for AI Agent (Claude Code):**
1. **JSON Output (`--json`)**: Every command that outputs data MUST support a `--json` flag.
2. **Exit Codes**: Always exit with `0` on success and non-zero (e.g., `1`) on failure.
3. **Core Delegation**: The CLI layer must contain NO business logic or raw SQL. It must call methods from `core.ProcessManager`.

## 2. Global Flags
- `--db string`: Path to the SQLite database file (default: `~/.taskctl/data.db`).
- `--json`: Output result in raw JSON format.

---

## 3. Command Reference

### Group A: Process Basic CRUD & Lifecycle

#### 1. `taskctl spawn` (Create Process)
- **Usage**: `taskctl spawn <title> [flags]`
- **Flags**:
  - `-d, --desc string`: Markdown description of the process.
  - `-p, --parent uint`: The PID of the parent process.
  - `--priority string`: `low`, `medium` (default), or `high`.
- **Core API**: `core.Spawn()`

#### 2. `taskctl ps` (Read / List Processes)
- **Usage**: `taskctl ps [flags]`
- **Flags**:
  - `-s, --status string`: Filter by status (default: `running`). Use `all` to show everything.
  - `-t, --tree`: Display output as a hierarchical process tree.
- **Core API**: `core.ListProcesses()` or `core.GetProcessTree()`

#### 3. `taskctl inspect` (Read Single Process) 🆕
- **Usage**: `taskctl inspect <pid>`
- **Description**: Shows the detailed metadata of a single process (Title, Description, Status, Priority, Timestamps). Similar to `docker inspect`.
- **Core API**: `core.GetProcess(pid)`

#### 4. `taskctl update` (Update Process Attributes) 🆕
- **Usage**: `taskctl update <pid> [flags]`
- **Description**: Modify basic attributes of a process without changing its running state. Useful for fixing typos or re-prioritizing.
- **Flags**:
  - `-t, --title string`: New title.
  - `-d, --desc string`: New description.
  - `--priority string`: Change priority.
  - `--ranking float`: Change ranking weight.
- **Core API**: `core.UpdateProcess(pid, updates)`

#### 5. `taskctl kill` (Delete Process)
- **Usage**: `taskctl kill <pid>`
- **Description**: Hard deletes the process and CASCADE deletes all its associated logs.
- **Core API**: `core.DeleteProcess(pid)`

---

### Group B: Logs Basic CRUD

#### 6. `taskctl log` (Create / Append Log)
- **Usage**: `taskctl log <pid> <content>`
- **Description**: Appends a free-form Markdown log to a process.
- **Core API**: `core.RecordLog()`

#### 7. `taskctl logs` (Read Logs of a Process) 🆕
- **Usage**: `taskctl logs <pid> [flags]`
- **Description**: Lists all log entries for a specific process in chronological order. Similar to `docker logs`.
- **Flags**:
  - `-n, --tail int`: Number of most recent logs to show (default: all).
- **Core API**: `core.GetProcessLogs(pid)`

#### 8. `taskctl log update` (Update/Edit Log) 🆕
- **Usage**: `taskctl log update <log_id> <new_content>`
- **Description**: Fix typos or update the content of a specific log entry.
- **Core API**: `core.UpdateLog(logID, newContent)`

#### 9. `taskctl log rm` (Delete Log) 🆕
- **Usage**: `taskctl log rm <log_id>`
- **Description**: Permanently deletes a specific log entry.
- **Core API**: `core.DeleteLog(logID)`

---

### Group C: State Machine (Workflow)

#### 10. `taskctl block` (Suspend)
- **Usage**: `taskctl block <pid> -m <reason>`
- **Description**: Changes status to `blocked` and automatically records a log with the reason.
- **Core API**: `core.ChangeState(pid, "blocked", reason)`

#### 11. `taskctl wake` (Resume)
- **Usage**: `taskctl wake <pid> -m <reason>`
- **Description**: Changes status back to `running` and logs the action.
- **Core API**: `core.ChangeState(pid, "running", reason)`

#### 12. `taskctl terminate` (Finish)
- **Usage**: `taskctl terminate <pid> -m <reason>`
- **Description**: Marks the process as `terminated` (done).
- **Core API**: `core.ChangeState(pid, "terminated", reason)`

---

### Group D: Analysis & Search

#### 13. `taskctl timeline` (Global Log Stream)
- **Usage**: `taskctl timeline [flags]`
- **Description**: Shows a chronological stream of logs across *all* processes.

#### 14. `taskctl grep` (Global Full-Text Search)
- **Usage**: `taskctl grep <keyword>`
- **Description**: Uses SQLite FTS5 to search across all titles, descriptions, and log contents.

#### 15. `taskctl stats` (Activity Heatmap)
- **Usage**: `taskctl stats`
- **Description**: Shows log counts per day for heatmaps.

---

### Group E: Export & UI

#### 16. `taskctl export` (Dump to Markdown)
- **Usage**: `taskctl export <pid> -o <dir>`
- **Description**: Generates a `.md` file containing process info and logs.

#### 17. `taskctl tui`
- **Usage**: `taskctl tui`
- **Description**: Launches the Bubble Tea interactive terminal application.

#### 18. `taskctl web`
- **Usage**: `taskctl web --port 8080`
- **Description**: Starts the embedded Vue web application and REST API server.

***

### 这次补充的亮点解析：

1. **`inspect` 命令**：完美契合 Unix/Docker 的习惯。`ps` 只能看个大概的列表，当你想要查看某个任务的完整描述（Description）和创建时间时，敲一下 `taskctl inspect 42`，直接打印出一个排版漂亮的详情卡片。
2. **`update` 拆分**：严格区分了“修改属性（`update`）”和“修改状态（`block/wake`）”。修改属性不会产生系统级的自动日志（因为只是改个错别字或者改个标题），而修改状态会强制产生日志，这保证了时间线的纯净。
3. **独立的 `logs` 命令**：之前只有全局视角的 `timeline`，现在补充了针对单个任务的 `taskctl logs <pid>`，这在你想要专门复盘某个任务的进度时极其方便。
4. **Log 的编辑与删除 (`log update`, `log rm`)**：对于个人工具，容错性很重要，手抖打错字或者挂错 PID 时，可以轻松撤回。

现在，这套 CLI 已经完全闭环了！无论是 AI 还是人类极客，使用起来都会感觉像在使用原生操作系统一样流畅。
