package tui

import (
	"fmt"
	"strings"

	"taskctl/core"

	"github.com/charmbracelet/lipgloss"
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
	if m.viewMode == ViewParentSelect {
		return m.parentSelectView()
	}
	if m.viewMode == ViewDeleteConfirm {
		return m.deleteConfirmView()
	}
	if m.viewMode == ViewExportConfirm {
		return m.exportConfirmView()
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
	var iconStyle lipgloss.Style
	var statusIcon string
	switch process.Status {
	case core.StatusRunning:
		iconStyle = statusStyle
		statusIcon = "▶"
	case core.StatusBlocked:
		iconStyle = blockedStyle
		statusIcon = "⏸"
	case core.StatusSuspended:
		iconStyle = suspendedStyle
		statusIcon = "⏹"
	case core.StatusTerminated:
		iconStyle = doneStyle
		statusIcon = "✓"
	default:
		iconStyle = helpStyle
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
		iconStyle.Render(statusIcon),
		priorityStyle.Render(priorityIcon),
		process.Title,
	)

	if idx == m.cursor && process.Description != "" {
		// Render description with markdown support if enabled
		var descContent string
		if m.markdownEnabled {
			// Use markdown rendering with limited width for list view
			mdRendered := RenderMarkdown(process.Description, 42)
			// Add proper indentation to each line
			mdLines := strings.Split(mdRendered, "\n")
			for i, mdLine := range mdLines {
				if mdLine == "" {
					continue
				}
				if i == 0 {
					descContent = "    └─ " + mdLine
				} else {
					descContent += "\n      " + mdLine
				}
			}
		} else {
			// Fallback to plain text with basic formatting
			desc := process.Description
			lines := strings.Split(desc, "\n")
			for i, descLine := range lines {
				if len(descLine) > 40 {
					descLine = descLine[:37] + "..."
				}
				if i == 0 {
					descContent = "    └─ " + descLine
				} else {
					descContent += "\n      " + descLine
				}
			}
		}
		line += "\n" + descContent
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
		filterStr = "等待中"
	} else if m.statusFilter == core.StatusTerminated {
		filterStr = "已终止"
	}
	markdownStatus := "开"
	if !m.markdownEnabled {
		markdownStatus = "关"
	}
	return helpStyle.Render(fmt.Sprintf(" 过滤:%s | MD:%s | j/k:导航  enter:详情  s:新建  d:删除  A:全部  /:搜索  q:退出  m:MD  ?:帮助",
		filterStr, markdownStatus)) + "\n" +
		helpStyle.Render(" r/B/S/T:过滤  G:时间线  D:统计  Y:树")
}

func (m Model) detailView() string {
	if m.currentProcess == nil {
		return "加载中..."
	}

	var b strings.Builder

	// Header
	b.WriteString(titleStyle.Render(fmt.Sprintf(" 进程 #%d: %s ", m.currentProcess.ID, m.currentProcess.Title)) + "\n")
	b.WriteString(borderStyle.Render(strings.Repeat("─", 50)) + "\n\n")

	// Show export success message if available
	if m.exportSuccessMsg != "" {
		b.WriteString(statusStyle.Render("✓ " + m.exportSuccessMsg) + "\n\n")
	}

	// Process info
	// Status in Chinese
	statusText := ""
	switch m.currentProcess.Status {
	case core.StatusRunning:
		statusText = "运行中"
	case core.StatusBlocked:
		statusText = "已阻塞"
	case core.StatusSuspended:
		statusText = "等待中"
	case core.StatusTerminated:
		statusText = "已终止"
	default:
		statusText = string(m.currentProcess.Status)
	}
	b.WriteString(fmt.Sprintf("状态:      %s %s\n",
		m.getStatusStyle(m.currentProcess.Status).Render(m.getStatusIcon(m.currentProcess.Status)),
		statusText))

	// Priority in Chinese
	priorityText := ""
	switch m.currentProcess.Priority {
	case core.PriorityHigh:
		priorityText = "高"
	case core.PriorityMedium:
		priorityText = "中"
	case core.PriorityLow:
		priorityText = "低"
	default:
		priorityText = string(m.currentProcess.Priority)
	}
	b.WriteString(fmt.Sprintf("优先级:    %s %s\n",
		m.getPriorityStyle(m.currentProcess.Priority).Render(m.getPriorityIcon(m.currentProcess.Priority)),
		priorityText))

	b.WriteString(fmt.Sprintf("创建时间:  %s\n",
		helpStyle.Render(m.currentProcess.CreatedAt.Format("2006-01-02 15:04"))))

	b.WriteString(fmt.Sprintf("更新时间:  %s\n",
		helpStyle.Render(m.currentProcess.UpdatedAt.Format("2006-01-02 15:04"))))

	// Parent process
	if m.currentProcess.ParentID != nil {
		parentTitle := "未知"
		if m.currentProcess.Parent != nil && m.currentProcess.Parent.Title != "" {
			parentTitle = m.currentProcess.Parent.Title
		}
		b.WriteString(fmt.Sprintf("父进程:    #%d %s\n",
			*m.currentProcess.ParentID,
			helpStyle.Render(parentTitle)))
	}

	if m.currentProcess.Description != "" {
		b.WriteString("\n描述:\n")
		var descContent string
		if m.markdownEnabled {
			descContent = RenderMarkdown(m.currentProcess.Description, 46)
		} else {
			descContent = "  " + m.currentProcess.Description
		}
		b.WriteString(descContent + "\n")
	}

	// Logs
	b.WriteString("\n" + borderStyle.Render(strings.Repeat("─", 50)) + "\n")
	b.WriteString(titleStyle.Render(" 时间线 ") + "\n")

	if len(m.processLogs) == 0 {
		b.WriteString(helpStyle.Render("暂无日志记录。") + "\n")
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

			// Render log content with or without markdown
			var logContent string
			if m.markdownEnabled {
				logContent = RenderMarkdown(log.Content, 35)
			} else {
				logContent = log.Content
			}

			b.WriteString(fmt.Sprintf(" %s [%s] %s %s\n",
				cursor,
				helpStyle.Render(timestamp),
				logStyle.Render(icon),
				logContent))
		}
	}

	b.WriteString("\n" + borderStyle.Render(strings.Repeat("─", 50)) + "\n")

	// Status bar
	b.WriteString(m.renderDetailStatusBar())

	return b.String()
}

func (m Model) renderDetailStatusBar() string {
	markdownStatus := "开"
	if !m.markdownEnabled {
		markdownStatus = "关"
	}
	return helpStyle.Render(" E:编辑  b/p/w/t:状态  a:添加日志  e:编辑日志  x:删除  X:导出  j/k:选择  q:返回") + "\n" +
		helpStyle.Render(fmt.Sprintf(" m:Markdown(%s)", markdownStatus))
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

	// Show different help text for state change (allows empty input)
	if m.pendingStateChange != nil {
		b.WriteString(helpStyle.Render(" enter:确认 (备注可选)  esc:取消"))
	} else if m.inputPrompt == "Search" {
		b.WriteString(helpStyle.Render(" enter:搜索  esc:取消"))
	} else {
		// Log input - multi-line textarea, show Ctrl+J as macOS alternative
		b.WriteString(helpStyle.Render(" Ctrl+J:确认 (macOS)  esc:取消"))
	}

	return b.String()
}

func (m Model) helpView() string {
	// All help content as lines
	lines := []string{
		titleStyle.Render(" 帮助 "),
		borderStyle.Render(strings.Repeat("─", 50)),
		"",
		titleStyle.Render("列表视图快捷键:"),
		"  j/k or ↑/↓    导航进程列表",
		"  1-9           快速跳转到第N项",
		"  enter         查看进程详情",
		"  s             新建进程",
		"  d             删除进程 (需确认)",
		"  A             显示全部进程",
		"  r/B/S/T       按状态过滤 (运行/阻塞/等待/终止)",
		"  /             全文搜索",
		"  G             全局时间线",
		"  D             活动统计",
		"  Y             进程树",
		"  ?             显示帮助",
		"  q/ctrl+c      退出程序",
		"",
		titleStyle.Render("详情视图快捷键:"),
		"  E             编辑进程信息 (标题/描述/优先级/父进程)",
		"  b             阻塞进程 (⏸) - 可选备注",
		"  p             等待中 (⏹) - 可选备注",
		"  w             唤醒进程 (▶) - 可选备注",
		"  t             终止进程 (✓) - 可选备注",
		"  a             添加日志",
		"  e             编辑选中的日志",
		"  x             删除选中的日志 (需确认)",
		"  X             导出为 Markdown 文件",
		"  m             切换 Markdown 渲染",
		"  j/k           选择日志",
		"  q/esc/h/←     返回列表",
		"",
		titleStyle.Render("创建/编辑进程快捷键:"),
		"  tab/shift+tab 切换字段",
		"  Ctrl+J        确认创建/编辑 (macOS)",
		"  enter         选择父进程",
		"  esc           取消",
		"",
		titleStyle.Render("选择父进程快捷键:"),
		"  j/k           导航进程列表",
		"  enter         选择父进程",
		"  esc/q         取消",
		"",
		titleStyle.Render("搜索视图快捷键:"),
		"  j/k           导航搜索结果",
		"  enter         跳转到进程详情",
		"  /             新搜索",
		"  q/esc         返回列表",
		"",
		titleStyle.Render("时间线视图快捷键:"),
		"  j/k           导航时间线",
		"  enter         跳转到进程详情",
		"  q/esc         返回列表",
		"",
		titleStyle.Render("统计视图快捷键:"),
		"  d             30天统计",
		"  h/H           7天/90天统计",
		"  D             365天统计",
		"  q/esc         返回列表",
		"",
		titleStyle.Render("树视图快捷键:"),
		"  j/k           导航进程树",
		"  enter         查看进程详情",
		"  space         展开/折叠",
		"  q/esc         返回列表",
		"",
		titleStyle.Render("状态图标:"),
		"  ▶             运行中 (running)",
		"  ⏸             已阻塞 (blocked)",
		"  ⏹             等待中 (suspended)",
		"  ✓             已终止 (terminated)",
		"",
		titleStyle.Render("优先级图标:"),
		"  H             高优先级",
		"  M             中优先级",
		"  L             低优先级",
		"",
		titleStyle.Render("Markdown 语法支持:"),
		"  **bold**      粗体文本",
		"  *italic*      斜体文本",
		"  `code`        行内代码",
		"  ```block```   代码块",
		"  # Header      一级标题",
		"  ## Header     二级标题",
		"  - item        列表项",
		"  > quote       引用文本",
		"  [text](url)   链接",
		"",
		borderStyle.Render(strings.Repeat("─", 50)),
		helpStyle.Render(" j/k:滚动  esc/?:关闭帮助"),
	}

	// Calculate visible lines (reserve 2 lines for top/bottom margins)
	maxLines := 30
	visibleLines := maxLines

	// Adjust offset if out of bounds
	if m.helpOffset < 0 {
		m.helpOffset = 0
	} else if m.helpOffset > len(lines)-visibleLines && len(lines) > visibleLines {
		m.helpOffset = len(lines) - visibleLines
	}

	// Calculate visible range
	start := m.helpOffset
	end := start + visibleLines
	if end > len(lines) {
		end = len(lines)
	}

	// Build visible output
	var b strings.Builder
	for i := start; i < end; i++ {
		b.WriteString(lines[i] + "\n")
	}

	// Add scroll indicator if needed
	if len(lines) > visibleLines {
		scrollHint := ""
		if m.helpOffset > 0 && end < len(lines) {
			scrollHint = helpStyle.Render(" ▲ 更多内容上下滚动 ▼ ")
		} else if m.helpOffset > 0 {
			scrollHint = helpStyle.Render(" ▲ 顶部 ")
		} else if end < len(lines) {
			scrollHint = helpStyle.Render(" ▼ 更多内容下方 ")
		}
		if scrollHint != "" {
			// Clear line and show hint at bottom
			b.WriteString("\r" + strings.Repeat(" ", 50) + "\r" + scrollHint + "\n")
		}
	}

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
	b.WriteString(descLabel + "\n" + m.spawnDesc.View() + "\n\n")

	// Priority field
	priorityLabel := "优先级:"
	if m.spawnFocusedField == 2 {
		priorityLabel = cursorStyle.Render(priorityLabel)
	}
	b.WriteString(priorityLabel + " " + m.spawnPriority.View() + " (H/M/L)\n\n")

	// Parent field
	parentLabel := "父进程:"
	if m.spawnFocusedField == 3 {
		parentLabel = cursorStyle.Render(parentLabel)
	}
	parentStr := "无"
	if m.selectedParentID != nil {
		if m.selectedParentName != "" {
			// Use cached name
			parentStr = fmt.Sprintf("#%d %s", *m.selectedParentID, m.selectedParentName)
		} else {
			// Fallback to showing just ID if name not cached
			parentStr = fmt.Sprintf("#%d", *m.selectedParentID)
		}
	}
	b.WriteString(parentLabel + " " + helpStyle.Render(parentStr) + "\n\n")

	b.WriteString(borderStyle.Render(strings.Repeat("─", 50)) + "\n")

	action := "创建"
	if m.editingProcessID > 0 {
		action = "保存"
	}
	b.WriteString(helpStyle.Render(fmt.Sprintf(" tab:切换字段  Ctrl+J:%s (macOS)  esc:取消", action)))

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

func (m Model) parentSelectView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(" 选择父进程 ") + "\n")
	b.WriteString(borderStyle.Render(strings.Repeat("─", 50)) + "\n\n")

	if len(m.availableParents) == 0 {
		b.WriteString(helpStyle.Render("没有可用的父进程。") + "\n")
	} else {
		// Show "None" option first
		cursor := " "
		if m.parentCursor == -1 {
			cursor = cursorStyle.Render("►")
		}
		selected := ""
		if m.selectedParentID == nil {
			selected = " [✓]"
		}
		b.WriteString(fmt.Sprintf("%s [无] %s\n\n", cursor, selected))

		// Show available parents
		for i, parent := range m.availableParents {
			cursor := " "
			if i == m.parentCursor {
				cursor = cursorStyle.Render("►")
			}

			selected := ""
			if m.selectedParentID != nil && parent.ID == *m.selectedParentID {
				selected = " [✓]"
			}

			// Show status icon
			statusIcon := "▶"
			if parent.Status == core.StatusBlocked {
				statusIcon = "⏸"
			} else if parent.Status == core.StatusSuspended {
				statusIcon = "⏹"
			} else if parent.Status == core.StatusTerminated {
				statusIcon = "✓"
			}

			b.WriteString(fmt.Sprintf("%s [#%d] %s %s%s\n",
				cursor, parent.ID, statusIcon, parent.Title, selected))
		}
	}

	b.WriteString("\n" + borderStyle.Render(strings.Repeat("─", 50)) + "\n")
	b.WriteString(helpStyle.Render(" j/k:导航  enter:选择  esc:取消"))

	return b.String()
}

func (m Model) deleteConfirmView() string {
	var b strings.Builder

	var title, itemType, itemName, warning string
	if m.confirmDeleteType == "process" {
		title = " 确认删除进程 "
		itemType = "进程"
		itemName = m.confirmDeleteName
		warning = "⚠️  警告：删除进程将同时删除其所有子进程和日志！"
	} else {
		title = " 确认删除日志 "
		itemType = "日志"
		itemName = m.confirmDeleteName
		warning = "⚠️  此操作无法撤销！"
	}

	b.WriteString(titleStyle.Render(title) + "\n")
	b.WriteString(borderStyle.Render(strings.Repeat("─", 50)) + "\n\n")

	b.WriteString(warningStyle.Render(warning) + "\n\n")
	b.WriteString(fmt.Sprintf("确定要删除此%s吗？\n\n", itemType))

	// Truncate name if too long
	displayName := itemName
	if len(displayName) > 40 {
		displayName = displayName[:40] + "..."
	}

	b.WriteString(helpStyle.Render(fmt.Sprintf("  %s: %s\n", itemType, displayName)))
	b.WriteString("\n")

	b.WriteString(borderStyle.Render(strings.Repeat("─", 50)) + "\n")
	b.WriteString(helpStyle.Render(" y:确认删除  n/esc:取消"))

	return b.String()
}

func (m Model) exportConfirmView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(" 导出进程 ") + "\n")
	b.WriteString(borderStyle.Render(strings.Repeat("─", 50)) + "\n\n")

	b.WriteString(fmt.Sprintf("确定要导出进程为 Markdown 文件吗？\n\n"))
	b.WriteString(helpStyle.Render(fmt.Sprintf("  文件名: %s\n", m.exportFileName)))
	b.WriteString(helpStyle.Render(fmt.Sprintf("  路径:   %s\n", m.exportFilePath)))
	b.WriteString("\n")

	b.WriteString(borderStyle.Render(strings.Repeat("─", 50)) + "\n")
	b.WriteString(helpStyle.Render(" y:确认导出  n/esc:取消"))

	return b.String()
}
