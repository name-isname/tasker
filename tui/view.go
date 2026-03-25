package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"taskctl/core"
)

// View renders the TUI
func (m Model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	if m.err != nil {
		return m.errorView()
	}

	if m.viewMode == ViewList {
		return m.listView()
	}
	if m.viewMode == ViewDetail {
		return m.detailView()
	}
	if m.viewMode == ViewInput {
		return m.inputView()
	}
	return m.helpView()
}

func (m Model) listView() string {
	var b strings.Builder

	// Header
	b.WriteString(titleStyle.Render(" taskctl - Process Manager ") + "\n")
	b.WriteString(borderStyle.Render(strings.Repeat("─", 50)) + "\n\n")

	// Process list
	if len(m.processes) == 0 {
		b.WriteString(helpStyle.Render("No processes found. Press 's' to spawn one.") + "\n")
	} else {
		// Calculate visible items based on viewport
		visibleCount := len(m.processes) - m.viewportOffset
		if visibleCount > 15 {
			visibleCount = 15
		}

		for i := 0; i < visibleCount; i++ {
			idx := m.viewportOffset + i
			if idx >= len(m.processes) {
				break
			}
			process := m.processes[idx]
			b.WriteString(m.renderProcessItem(idx, process))
		}
	}

	b.WriteString("\n")
	b.WriteString(borderStyle.Render(strings.Repeat("─", 50)) + "\n")

	// Status bar
	b.WriteString(m.renderStatusBar())

	return b.String()
}

func (m Model) renderProcessItem(idx int, process core.Process) string {
	cursor := " "
	if idx == m.cursor {
		cursor = "►"
	}

	// Status icon and style
	var statusStyle lipgloss.Style
	var statusIcon string
	switch process.Status {
	case core.StatusRunning:
		statusStyle = statusStyle
		statusIcon = "▶"
	case core.StatusBlocked:
		statusStyle = blockedStyle
		statusIcon = "⏸"
	case core.StatusSuspended:
		statusStyle = suspendedStyle
		statusIcon = "⏹"
	case core.StatusTerminated:
		statusStyle = doneStyle
		statusIcon = "✓"
	default:
		statusStyle = helpStyle
		statusIcon = "?"
	}

	// Priority icon
	var priorityStyle lipgloss.Style
	var priorityIcon string
	switch process.Priority {
	case core.PriorityHigh:
		priorityStyle = priorityHighStyle
		priorityIcon = "H"
	case core.PriorityMedium:
		priorityStyle = priorityMediumStyle
		priorityIcon = "M"
	case core.PriorityLow:
		priorityStyle = priorityLowStyle
		priorityIcon = "L"
	default:
		priorityStyle = helpStyle
		priorityIcon = "?"
	}

	// Build the line
	line := fmt.Sprintf("%s [%s] [%s] %s",
		cursorStyle.Render(cursor),
		statusStyle.Render(statusIcon),
		priorityStyle.Render(priorityIcon),
		process.Title,
	)

	if idx == m.cursor && process.Description != "" {
		// Truncate description if too long
		desc := process.Description
		if len(desc) > 40 {
			desc = desc[:37] + "..."
		}
		line += "\n" + helpStyle.Render("    └─ "+desc)
	}

	return line + "\n"
}

func (m Model) renderStatusBar() string {
	return helpStyle.Render(" j/k:导航  enter:详情  s:新建  q:退出  ?:帮助")
}

func (m Model) detailView() string {
	if m.currentProcess == nil {
		return "Loading..."
	}

	var b strings.Builder

	// Header
	b.WriteString(titleStyle.Render(fmt.Sprintf(" Process #%d: %s ", m.currentProcess.ID, m.currentProcess.Title)) + "\n")
	b.WriteString(borderStyle.Render(strings.Repeat("─", 50)) + "\n\n")

	// Process info
	b.WriteString(fmt.Sprintf("Status:    %s %s\n",
		m.getStatusStyle(m.currentProcess.Status).Render(m.getStatusIcon(m.currentProcess.Status)),
		m.currentProcess.Status))

	b.WriteString(fmt.Sprintf("Priority:  %s %s\n",
		m.getPriorityStyle(m.currentProcess.Priority).Render(m.getPriorityIcon(m.currentProcess.Priority)),
		m.currentProcess.Priority))

	b.WriteString(fmt.Sprintf("Created:   %s\n",
		helpStyle.Render(m.currentProcess.CreatedAt.Format("2006-01-02 15:04"))))

	b.WriteString(fmt.Sprintf("Updated:   %s\n",
		helpStyle.Render(m.currentProcess.UpdatedAt.Format("2006-01-02 15:04"))))

	if m.currentProcess.Description != "" {
		b.WriteString("\nDescription:\n")
		b.WriteString(helpStyle.Render("  "+m.currentProcess.Description) + "\n")
	}

	// Logs
	b.WriteString("\n" + borderStyle.Render(strings.Repeat("─", 50)) + "\n")
	b.WriteString(titleStyle.Render(" Timeline ") + "\n")

	if len(m.processLogs) == 0 {
		b.WriteString(helpStyle.Render("No logs yet.") + "\n")
	} else {
		for _, log := range m.processLogs {
			timestamp := log.CreatedAt.Format("15:04")
			icon := "📝"
			logStyle := logProgressStyle
			if log.LogType == core.LogTypeStateChange {
				icon = "🔄"
				logStyle = logStateChangeStyle
			}
			b.WriteString(fmt.Sprintf("   [%s] %s %s\n",
				helpStyle.Render(timestamp),
				logStyle.Render(icon),
				log.Content))
		}
	}

	b.WriteString("\n" + borderStyle.Render(strings.Repeat("─", 50)) + "\n")

	// Status bar
	b.WriteString(m.renderDetailStatusBar())

	return b.String()
}

func (m Model) renderDetailStatusBar() string {
	return helpStyle.Render(" b:block  w:wake  t:terminate  a:add log  esc:back  q:quit")
}

func (m Model) errorView() string {
	return fmt.Sprintf("Error: %v\n\nPress q to quit", m.err)
}

// Helper methods for styling
func (m Model) getStatusIcon(status core.ProcessStatus) string {
	switch status {
	case core.StatusRunning:
		return "▶"
	case core.StatusBlocked:
		return "⏸"
	case core.StatusSuspended:
		return "⏹"
	case core.StatusTerminated:
		return "✓"
	default:
		return "?"
	}
}

func (m Model) getStatusStyle(status core.ProcessStatus) lipgloss.Style {
	switch status {
	case core.StatusRunning:
		return statusStyle
	case core.StatusBlocked:
		return blockedStyle
	case core.StatusSuspended:
		return suspendedStyle
	case core.StatusTerminated:
		return doneStyle
	default:
		return helpStyle
	}
}

func (m Model) getPriorityIcon(priority core.ProcessPriority) string {
	switch priority {
	case core.PriorityHigh:
		return "H"
	case core.PriorityMedium:
		return "M"
	case core.PriorityLow:
		return "L"
	default:
		return "?"
	}
}

func (m Model) getPriorityStyle(priority core.ProcessPriority) lipgloss.Style {
	switch priority {
	case core.PriorityHigh:
		return priorityHighStyle
	case core.PriorityMedium:
		return priorityMediumStyle
	case core.PriorityLow:
		return priorityLowStyle
	default:
		return helpStyle
	}
}

func (m Model) inputView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(" "+m.inputPrompt+" ") + "\n")
	b.WriteString(borderStyle.Render(strings.Repeat("─", 50)) + "\n\n")

	b.WriteString(m.textInput.View())
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render(" enter:确认  esc:取消"))

	return b.String()
}

func (m Model) helpView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(" 帮助 ") + "\n")
	b.WriteString(borderStyle.Render(strings.Repeat("─", 50)) + "\n\n")

	b.WriteString(titleStyle.Render("列表视图快捷键:\n"))
	b.WriteString("  j/k or ↑/↓    导航进程列表\n")
	b.WriteString("  1-9           快速跳转到第N项\n")
	b.WriteString("  enter         查看进程详情\n")
	b.WriteString("  s             新建进程 (TODO)\n")
	b.WriteString("  ?             显示帮助\n")
	b.WriteString("  q/ctrl+c      退出\n\n")

	b.WriteString(titleStyle.Render("详情视图快捷键:\n"))
	b.WriteString("  b             阻塞进程 (⏸)\n")
	b.WriteString("  w             唤醒进程 (▶)\n")
	b.WriteString("  t             终止进程 (✓)\n")
	b.WriteString("  a             添加日志\n")
	b.WriteString("  esc/h/←       返回列表\n")
	b.WriteString("  q             退出\n\n")

	b.WriteString(titleStyle.Render("状态图标:\n"))
	b.WriteString("  ▶             运行中 (running)\n")
	b.WriteString("  ⏸             已阻塞 (blocked)\n")
	b.WriteString("  ⏹             已暂停 (suspended)\n")
	b.WriteString("  ✓             已终止 (terminated)\n\n")

	b.WriteString(titleStyle.Render("优先级图标:\n"))
	b.WriteString("  H             高优先级\n")
	b.WriteString("  M             中优先级\n")
	b.WriteString("  L             低优先级\n\n")

	b.WriteString(borderStyle.Render(strings.Repeat("─", 50)) + "\n")
	b.WriteString(helpStyle.Render(" esc/?:关闭帮助"))

	return b.String()
}
