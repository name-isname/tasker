# taskctl

> 三合一任务管理工具（CLI + TUI + Web）

## 概述

**taskctl** 是一个用 Go 编写的任务管理工具，通过单个二进制文件支持三种不同的界面：

- **CLI** - 命令行界面，快速操作任务
- **TUI** - 交互式终端用户界面
- **Web** - 浏览器界面

### 架构设计

项目采用 **"核心库优先"（Core-Library First）** 设计模式：

```
┌─────────────────────────────────────────────────────┐
│                    core/                            │
│            (业务逻辑与数据层)                         │
│     AddTask() | ListTasks() | CompleteTask()        │
└──────────┬──────────────────────────────────────────┘
           │
     ┌─────┴─────┬──────────────┬─────────────┐
     ▼           ▼              ▼             ▼
  cli/        tui/            web/        (未来扩展)
(Cobra)    (Bubble Tea)      (Gin+Vue)
```

所有数据库和业务逻辑都位于 `core/` 包中。各表现层（CLI/TUI/Web）只需调用核心函数。

## 技术栈

| 层级 | 技术 |
|-----|-----|
| **数据库** | SQLite（使用 `github.com/glebarez/sqlite`，纯 Go 实现，无 CGO） |
| **ORM** | GORM |
| **CLI** | Cobra |
| **TUI** | Bubble Tea + Lipgloss |
| **Web API** | Gin |
| **前端** | Vue 3 + TypeScript + Vite + TailwindCSS |

## 安装

### 前置要求

- Go >= 1.22
- Node.js >= 18（前端开发需要）

### 从源码构建

```bash
# 克隆仓库
git clone https://github.com/yourusername/taskctl.git
cd taskctl

# 安装前端依赖
make install-frontend

# 构建所有内容（前端 + Go 二进制）
make build

# 生成的二进制文件为 ./taskctl
```

### 跨平台构建

```bash
make build-linux    # Linux AMD64
make build-mac      # macOS AMD64 + ARM64
make build-windows  # Windows AMD64
```

## 使用方法

### CLI 命令行

```bash
# 添加任务
./taskctl add "购买日用品"

# 列出所有任务
./taskctl list

# 以 JSON 格式输出（供 AI 代理使用）
./taskctl list --json

# 标记任务为完成
./taskctl complete 1

# 删除任务
./taskctl delete 1

# 指定数据库位置
./taskctl --db ~/mytasks.db list
```

### TUI 终端界面

```bash
# 启动交互式终端界面
./taskctl tui
```

使用 `j/k` 或 `↑/↓` 导航，`q` 退出。

### Web 界面

```bash
# 启动 Web 服务器（默认端口：8080）
./taskctl web

# 指定端口
./taskctl web --port 3000
```

然后访问 http://localhost:8080

## 开发

```bash
# 安装前端依赖
make install-frontend

# 启动前端开发服务器
make dev

# 直接运行 CLI（无需构建）
make run

# 仅构建前端
make build-frontend

# 仅构建 Go 二进制
make build-go

# 清理构建产物
make clean
```

## API 接口

Web 服务器在 `/api/v1/` 路径下提供 RESTful API：

- `GET /api/v1/tasks` - 获取所有任务
- `POST /api/v1/tasks` - 创建新任务
- `PUT /api/v1/tasks/:id/complete` - 标记任务完成
- `DELETE /api/v1/tasks/:id` - 删除任务

## 项目结构

```
taskctl/
├── main.go              # 入口文件
├── core/                # 业务逻辑与数据层
│   ├── db.go           # 数据库初始化
│   └── task.go         # Task 模型和 CRUD 操作
├── cli/                 # CLI 命令（Cobra）
│   ├── root.go         # 根命令
│   ├── add.go          # 添加任务
│   ├── list.go         # 列出任务（支持 --json）
│   ├── complete.go     # 完成任务
│   ├── delete.go       # 删除任务
│   ├── tui.go          # TUI 入口
│   └── web.go          # Web 服务器入口
├── tui/                 # 终端界面（Bubble Tea）
│   ├── model.go        # TUI 模型
│   └── view.go         # TUI 视图渲染
├── web/                 # Web 服务器和嵌入式前端
│   ├── server.go       # Gin API 服务器
│   ├── embed.go        # 静态文件嵌入
│   └── frontend/       # Vue + Vite + TailwindCSS
│       ├── src/
│       ├── dist/       # 构建产物（嵌入到二进制）
│       └── ...
├── Makefile            # 构建命令
├── CLAUDE.md           # AI 助手指南
├── README.md           # 英文文档
└── README.zh.md        # 中文文档
```

## 设计原则

1. **核心库优先**：`core/` 包包含所有业务逻辑和数据访问。表现层只负责用户交互。
2. **纯 Go SQLite**：使用 `glebarez/sqlite` 而非 `mattn/go-sqlite3`，避免 CGO 依赖，便于交叉编译。
3. **AI 友好**：CLI 支持 `--json` 标志，便于 AI 代理解析。
4. **单二进制分发**：前端资源通过 `//go:embed` 嵌入到 Go 二进制文件中。

## 许可证

MIT License
