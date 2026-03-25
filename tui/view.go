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
	if m.viewMode == ViewHelp {
		return m.helpView()
	}
	if m.viewMode == ViewSpawn || m.viewMode == ViewEditProcess {
		return m.spawnView()
	}
	if m.viewMode == ViewSearch {
		return m.searchView()
	}
	if m.viewMode == ViewTimeline {
		return m.timelineView()
	}
	if m.viewMode == ViewStats {
		return m.statsView()
	}
	if m.viewMode == ViewTree {
		return m.treeView()
	}
	return ""
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
	filterStr := "全部"
	if m.statusFilter == core.StatusRunning {
		filterStr = "运行中"
	} else if m.statusFilter == core.StatusBlocked {
		filterStr = "已阻塞"
	} else if m.statusFilter == core.StatusSuspended {
		filterStr = "已暂停"
	} else if m.statusFilter == core.StatusTerminated {
		filterStr = "已终止"
	}
	return helpStyle.Render(fmt.Sprintf(" 过滤:%s | j/k:导航  enter:详情  s:新建  r/B/S/T:过滤  /:搜索  G:时间线  D:统计  Y:树  q:返回/退出  ?:帮助", filterStr))
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
		for i, log := range m.processLogs {
			timestamp := log.CreatedAt.Format("15:04")
			icon := "📝"
			logStyle := logProgressStyle
			if log.LogType == core.LogTypeStateChange {
				icon = "🔄"
				logStyle = logStateChangeStyle
			}

			// Add cursor indicator for selected log
			cursor := " "
			if i == m.logCursor {
				cursor = cursorStyle.Render("►")
			} else {
				cursor = " "
			}

			b.WriteString(fmt.Sprintf(" %s [%s] %s %s\n",
				cursor,
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
	return helpStyle.Render(" E:编辑进程  b:阻塞  w:唤醒  t:终止  a:添加日志  e:编辑日志  j/k:选择日志  q/esc:返回")
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

	b.WriteString(titleStyle.Render("列表视图快捷键:") + "\n")
	b.WriteString("  j/k or ↑/↓    导航进程列表\n")
	b.WriteString("  1-9           快速跳转到第N项\n")
	b.WriteString("  enter         查看进程详情\n")
	b.WriteString("  s             新建进程\n")
	b.WriteString("  r/B/S/T/A     按状态过滤 (运行/阻塞/暂停/终止/全部)\n")
	b.WriteString("  /             全文搜索\n")
	b.WriteString("  G             全局时间线\n")
	b.WriteString("  D             活动统计\n")
	b.WriteString("  Y             进程树\n")
	b.WriteString("  ?             显示帮助\n")
	b.WriteString("  q/ctrl+c      退出程序\n\n")

	b.WriteString(titleStyle.Render("详情视图快捷键:") + "\n")
	b.WriteString("  E             编辑进程信息 (标题/描述/优先级)\n")
	b.WriteString("  b             阻塞进程 (⏸)\n")
	b.WriteString("  w             唤醒进程 (▶)\n")
	b.WriteString("  t             终止进程 (✓)\n")
	b.WriteString("  a             添加日志\n")
	b.WriteString("  e             编辑选中的日志\n")
	b.WriteString("  j/k           选择日志\n")
	b.WriteString("  q/esc/h/←     返回列表\n\n")

	b.WriteString(titleStyle.Render("搜索视图快捷键:") + "\n")
	b.WriteString("  j/k           导航搜索结果\n")
	b.WriteString("  enter         跳转到进程详情\n")
	b.WriteString("  /             新搜索\n")
	b.WriteString("  q/esc         返回列表\n\n")

	b.WriteString(titleStyle.Render("时间线视图快捷键:") + "\n")
	b.WriteString("  j/k           导航时间线\n")
	b.WriteString("  enter         跳转到进程详情\n")
	b.WriteString("  q/esc         返回列表\n\n")

	b.WriteString(titleStyle.Render("统计视图快捷键:") + "\n")
	b.WriteString("  d             30天统计\n")
	b.WriteString("  h/H           7天/90天统计\n")
	b.WriteString("  D             365天统计\n")
	b.WriteString("  q/esc         返回列表\n\n")

	b.WriteString(titleStyle.Render("树视图快捷键:") + "\n")
	b.WriteString("  j/k           导航进程树\n")
	b.WriteString("  enter         查看进程详情\n")
	b.WriteString("  space         展开/折叠\n")
	b.WriteString("  q/esc         返回列表\n\n")

	b.WriteString(titleStyle.Render("状态图标:") + "\n")
	b.WriteString("  ▶             运行中 (running)\n")
	b.WriteString("  ⏸             已阻塞 (blocked)\n")
	b.WriteString("  ⏹             已暂停 (suspended)\n")
	b.WriteString("  ✓             已终止 (terminated)\n\n")

	b.WriteString(titleStyle.Render("优先级图标:") + "\n")
	b.WriteString("  H             高优先级\n")
	b.WriteString("  M             中优先级\n")
	b.WriteString("  L             低优先级\n\n")

	b.WriteString(borderStyle.Render(strings.Repeat("─", 50)) + "\n")
	b.WriteString(helpStyle.Render(" esc/?:关闭帮助"))

	return b.String()
}

func (m Model) spawnView() string {
	var b strings.Builder

	// Show different title based on mode
	title := " 新建进程 "
	if m.editingProcessID > 0 {
		title = " 编辑进程 "
	}
	b.WriteString(titleStyle.Render(title) + "\n")
	b.WriteString(borderStyle.Render(strings.Repeat("─", 50)) + "\n\n")

	// Title field
	titleLabel := "标题:"
	if m.spawnFocusedField == 0 {
		titleLabel = cursorStyle.Render(titleLabel)
	}
	b.WriteString(titleLabel + " " + m.spawnTitle.View() + "\n\n")

	// Description field
	descLabel := "描述:"
	if m.spawnFocusedField == 1 {
		descLabel = cursorStyle.Render(descLabel)
	}
	b.WriteString(descLabel + " " + m.spawnDesc.View() + "\n\n")

	// Priority field
	priorityLabel := "优先级:"
	if m.spawnFocusedField == 2 {
		priorityLabel = cursorStyle.Render(priorityLabel)
	}
	b.WriteString(priorityLabel + " " + m.spawnPriority.View() + " (H/M/L)\n\n")

	b.WriteString(borderStyle.Render(strings.Repeat("─", 50)) + "\n")
	action := "创建"
	if m.editingProcessID > 0 {
		action = "保存"
	}
	b.WriteString(helpStyle.Render(fmt.Sprintf(" tab:切换字段  enter:%s  esc:取消", action)))

	return b.String()
}

func (m Model) searchView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(fmt.Sprintf(" 搜索: %s ", m.searchKeyword)) + "\n")
	b.WriteString(borderStyle.Render(strings.Repeat("─", 50)) + "\n\n")

	if len(m.searchResults) == 0 {
		b.WriteString(helpStyle.Render("未找到结果。") + "\n")
	} else {
		b.WriteString(fmt.Sprintf("找到 %d 个结果:\n\n", len(m.searchResults)))

		// Show visible results
		visibleCount := len(m.searchResults)
		if visibleCount > 15 {
			visibleCount = 15
		}

		for i := 0; i < visibleCount; i++ {
			result := m.searchResults[i]
			cursor := " "
			if i == m.searchCursor {
				cursor = cursorStyle.Render("►")
			}

			if result.Type == "process" {
				b.WriteString(fmt.Sprintf("%s [进程 #%d] %s\n",
					cursor, result.ID, result.Title))
				if result.Content != "" {
					preview := result.Content
					if len(preview) > 50 {
						preview = preview[:47] + "..."
					}
					b.WriteString(helpStyle.Render("    └─ "+preview) + "\n")
				}
			} else {
				icon := "📝"
				b.WriteString(fmt.Sprintf("%s [%s] 日志 #%d (进程 #%d)\n",
					cursor, icon, result.ID, result.ProcessID))
				content := result.Content
				if len(content) > 50 {
					content = content[:47] + "..."
				}
				b.WriteString(helpStyle.Render("    └─ "+content) + "\n")
			}
		}
	}

	b.WriteString("\n" + borderStyle.Render(strings.Repeat("─", 50)) + "\n")
	b.WriteString(helpStyle.Render(" j/k:导航  enter:查看  /:新搜索  q:返回"))

	return b.String()
}

func (m Model) timelineView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(" 全局时间线 ") + "\n")
	b.WriteString(borderStyle.Render(strings.Repeat("─", 50)) + "\n\n")

	if len(m.timelineEntries) == 0 {
		b.WriteString(helpStyle.Render("没有日志记录。") + "\n")
	} else {
		// Group by date
		currentDate := ""
		visibleCount := 0
		maxVisible := 20

		for i, entry := range m.timelineEntries {
			if visibleCount >= maxVisible {
				break
			}

			date := entry.CreatedAt.Format("2006-01-02")
			if date != currentDate {
				if currentDate != "" {
					b.WriteString("\n")
				}
				b.WriteString(titleStyle.Render("📅 "+date) + "\n")
				b.WriteString(borderStyle.Render(strings.Repeat("─", 50)) + "\n")
				currentDate = date
			}

			timeStr := entry.CreatedAt.Format("15:04")
			icon := "📝"
			if entry.LogType == core.LogTypeStateChange {
				icon = "🔄"
			}

			cursor := " "
			if i == m.timelineCursor {
				cursor = cursorStyle.Render("►")
			}

			// Truncate content if too long
			content := entry.Content
			if len(content) > 45 {
				content = content[:42] + "..."
			}

			b.WriteString(fmt.Sprintf(" %s [%s] %s %s: %s\n",
				cursor, timeStr, icon, entry.ProcessTitle, content))
			visibleCount++
		}
	}

	b.WriteString("\n" + borderStyle.Render(strings.Repeat("─", 50)) + "\n")
	b.WriteString(helpStyle.Render(" j/k:导航  enter:查看  q:返回"))

	return b.String()
}

func (m Model) statsView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(fmt.Sprintf(" 活动统计 (过去%d天) ", m.statsDays)) + "\n")
	b.WriteString(borderStyle.Render(strings.Repeat("─", 50)) + "\n\n")

	if len(m.statsData) == 0 {
		b.WriteString(helpStyle.Render("没有活动数据。") + "\n")
	} else {
		// Find max count for scaling
		maxCount := 0
		for _, stat := range m.statsData {
			if stat.Count > maxCount {
				maxCount = stat.Count
			}
		}

		// Print activity chart
		for _, stat := range m.statsData {
			// Create a simple bar chart
			barWidth := int(float64(stat.Count) / float64(maxCount) * 30)
			if barWidth == 0 && stat.Count > 0 {
				barWidth = 1
			}
			bar := strings.Repeat("█", barWidth)

			// Format date as MM-DD
			dateStr := stat.Date[5:] // Skip YYYY-

			b.WriteString(fmt.Sprintf(" %s │%s %d\n", dateStr, bar, stat.Count))
		}

		// Calculate totals
		totalLogs := 0
		for _, stat := range m.statsData {
			totalLogs += stat.Count
		}
		avgPerDay := float64(totalLogs) / float64(len(m.statsData))

		b.WriteString("\n" + borderStyle.Render(strings.Repeat("─", 50)) + "\n")
		b.WriteString(fmt.Sprintf(" 总计: %d 条日志 | 平均: %.1f 条/天\n", totalLogs, avgPerDay))
	}

	b.WriteString(borderStyle.Render(strings.Repeat("─", 50)) + "\n")
	b.WriteString(helpStyle.Render(" h:7天  d:30天  H:90天  D:365天  q:返回"))

	return b.String()
}

func (m Model) treeView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(" 进程树 ") + "\n")
	b.WriteString(borderStyle.Render(strings.Repeat("─", 50)) + "\n\n")

	nodes := m.flattenTreeNodes(m.treeNodes, 0)
	if len(nodes) == 0 {
		b.WriteString(helpStyle.Render("没有进程。") + "\n")
	} else {
		// Show visible nodes with viewport
		visibleCount := len(nodes)
		if visibleCount > 15 {
			visibleCount = 15
		}

		for i := 0; i < visibleCount; i++ {
			node := nodes[i]
			cursor := " "
			if i == m.treeCursor {
				cursor = cursorStyle.Render("►")
			}

			// Get status icon
			var statusIcon string
			switch node.Status {
			case core.StatusRunning:
				statusIcon = "🚀"
			case core.StatusBlocked:
				statusIcon = "⏸"
			case core.StatusSuspended:
				statusIcon = "⏹"
			case core.StatusTerminated:
				statusIcon = "✓"
			default:
				statusIcon = "?"
			}

			// Calculate depth for indentation
			depth := m.getNodeDepth(&node, m.treeNodes, 0)
			indent := strings.Repeat("  ", depth)

			// Show expand indicator if has children
			expandIndicator := ""
			if len(node.Children) > 0 {
				if m.treeExpanded[node.ID] {
					expandIndicator = "[-]"
				} else {
					expandIndicator = "[+]"
				}
			}

			b.WriteString(fmt.Sprintf("%s%s [#%d] %s %s %s\n",
				cursor, indent, node.ID, statusIcon, node.Title, expandIndicator))
		}
	}

	b.WriteString("\n" + borderStyle.Render(strings.Repeat("─", 50)) + "\n")
	b.WriteString(helpStyle.Render(" j/k:导航  enter:查看  space:展开/折叠  q:返回"))

	return b.String()
}

// getNodeDepth calculates the depth of a node in the tree
func (m Model) getNodeDepth(node *core.ProcessNode, nodes []*core.ProcessNode, depth int) int {
	for _, n := range nodes {
		if n.ID == node.ID {
			return depth
		}
		if len(n.Children) > 0 {
			if d := m.getNodeDepth(node, n.Children, depth+1); d >= 0 {
				return d
			}
		}
	}
	return -1
}
