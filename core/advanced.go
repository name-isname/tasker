package core

import (
	"fmt"
	"strings"
	"time"
	"gorm.io/gorm"
)

// ChangeProcessState atomically changes process status and creates a state-change log
// This is the preferred method for status transitions as it ensures data consistency
func ChangeProcessState(processID uint, newStatus ProcessStatus, reason string) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// Get current process
		var process Process
		if err := tx.First(&process, processID).Error; err != nil {
			return err
		}

		// Update status
		if err := tx.Model(&Process{}).Where("id = ?", processID).Update("status", newStatus).Error; err != nil {
			return err
		}

		// Create state change log
		logContent := fmt.Sprintf("Status changed from %s to %s", process.Status, newStatus)
		if reason != "" {
			logContent += fmt.Sprintf(". Reason: %s", reason)
		}

		log := &Log{
			ProcessID: processID,
			LogType:   LogTypeStateChange,
			Content:   logContent,
		}

		return tx.Create(log).Error
	})
}

// GetProcessTree retrieves a process and all its descendants as a tree structure
func GetProcessTree(rootID uint) (*ProcessNode, error) {
	var root Process
	if err := DB.First(&root, rootID).Error; err != nil {
		return nil, err
	}

	node := &ProcessNode{Process: root}
	if err := loadChildren(node); err != nil {
		return nil, err
	}

	return node, nil
}

// loadChildren recursively loads children for a process node
func loadChildren(node *ProcessNode) error {
	var children []Process
	if err := DB.Where("parent_id = ?", node.ID).Order("ranking DESC, created_at DESC").Find(&children).Error; err != nil {
		return err
	}

	node.Children = make([]*ProcessNode, len(children))
	for i := range children {
		childNode := &ProcessNode{Process: children[i]}
		if err := loadChildren(childNode); err != nil {
			return err
		}
		node.Children[i] = childNode
	}

	return nil
}

// GetFullProcessTree retrieves all root processes with their complete trees
func GetFullProcessTree() ([]*ProcessNode, error) {
	var roots []Process
	if err := DB.Where("parent_id IS NULL").Order("ranking DESC, created_at DESC").Find(&roots).Error; err != nil {
		return nil, err
	}

	nodes := make([]*ProcessNode, len(roots))
	for i := range roots {
		node := &ProcessNode{Process: roots[i]}
		if err := loadChildren(node); err != nil {
			return nil, err
		}
		nodes[i] = node
	}

	return nodes, nil
}

// GetTimeline retrieves logs across all processes within a time range
func GetTimeline(startTime, endTime time.Time, limit int) ([]TimelineEntry, error) {
	var entries []TimelineEntry

	query := DB.Table("logs").
		Select("logs.id, logs.process_id, logs.log_type, logs.content, logs.created_at, processes.title as process_title").
		Joins("JOIN processes ON processes.id = logs.process_id").
		Order("logs.created_at DESC")

	if !startTime.IsZero() {
		query = query.Where("logs.created_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		query = query.Where("logs.created_at <= ?", endTime)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Scan(&entries).Error
	return entries, err
}

// GetTodayTimeline retrieves all logs from today
func GetTodayTimeline() ([]TimelineEntry, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	return GetTimeline(startOfDay, endOfDay, 0)
}

// GetActivityStats retrieves daily log counts for the past N days
func GetActivityStats(days int) ([]ActivityStat, error) {
	// SQLite date format: YYYY-MM-DD
	sql := `
		SELECT DATE(created_at) as date, COUNT(*) as count
		FROM logs
		WHERE DATE(created_at) >= DATE('now', '-' || ? || ' days')
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`

	var stats []ActivityStat
	err := DB.Raw(sql, days).Scan(&stats).Error
	return stats, err
}

// GlobalSearch performs full-text search across processes and logs
func GlobalSearch(keyword string) ([]SearchResult, error) {
	if keyword == "" {
		return []SearchResult{}, nil
	}

	var results []SearchResult

	// Search in processes
	var processes []Process
	if err := DB.Where("title LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Find(&processes).Error; err != nil {
		return nil, err
	}

	for _, p := range processes {
		results = append(results, SearchResult{
			Type:    "process",
			ID:      p.ID,
			Title:   p.Title,
			Content: p.Description,
		})
	}

	// Search in logs
	var logs []Log
	if err := DB.Where("content LIKE ?", "%"+keyword+"%").Find(&logs).Error; err != nil {
		return nil, err
	}

	for _, l := range logs {
		results = append(results, SearchResult{
			Type:      "log",
			ID:        l.ID,
			Content:   l.Content,
			ProcessID: l.ProcessID,
		})
	}

	return results, nil
}

// GetProcessContext retrieves a process with all its logs for export
func GetProcessContext(processID uint) (*ProcessExport, error) {
	var process Process
	if err := DB.First(&process, processID).Error; err != nil {
		return nil, err
	}

	// Get all logs, ordered chronologically (oldest first for export)
	var logs []Log
	if err := DB.Where("process_id = ?", processID).Order("created_at ASC").Find(&logs).Error; err != nil {
		return nil, err
	}

	// Get child processes recursively
	var children []Process
	if err := DB.Where("parent_id = ?", processID).Find(&children).Error; err != nil {
		return nil, err
	}

	export := &ProcessExport{
		Process: process,
		Logs:    logs,
	}

	// Recursively get children context
	for _, child := range children {
		childExport, err := GetProcessContext(child.ID)
		if err != nil {
			return nil, err
		}
		export.Children = append(export.Children, *childExport)
	}

	return export, nil
}

// FormatProcessTree renders a process tree as ASCII art (like pstree)
func FormatProcessTree(node *ProcessNode, prefix string, isLast bool) string {
	var sb strings.Builder

	// Process icon based on status
	icon := getProcessIcon(node.Status)

	// Connector
	connector := "├── "
	if isLast {
		connector = "└── "
	}

	sb.WriteString(fmt.Sprintf("%s%s[#%d] %s %s", prefix, connector, node.ID, icon, node.Title))

	// Add status details
	if node.Status == StatusBlocked {
		// Get latest log for reason
		var latestLog Log
		if err := DB.Where("process_id = ? AND log_type = ?", node.ID, LogTypeStateChange).
			Order("created_at DESC").First(&latestLog).Error; err == nil {
			if strings.Contains(latestLog.Content, "Reason:") {
				reason := strings.Split(latestLog.Content, "Reason: ")[1]
				sb.WriteString(fmt.Sprintf(" (%s)", reason))
			}
		}
	}

	sb.WriteString("\n")

	// Process children
	for i, child := range node.Children {
		isLastChild := i == len(node.Children)-1
		newPrefix := prefix
		if isLast {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		sb.WriteString(FormatProcessTree(child, newPrefix, isLastChild))
	}

	return sb.String()
}

// FormatFullTree renders the entire process forest
func FormatFullTree() (string, error) {
	nodes, err := GetFullProcessTree()
	if err != nil {
		return "", err
	}

	if len(nodes) == 0 {
		return "No processes found.\n", nil
	}

	var sb strings.Builder
	for i, node := range nodes {
		isLast := i == len(nodes)-1
		sb.WriteString(FormatProcessTree(node, "", isLast))
	}

	return sb.String(), nil
}

func getProcessIcon(status ProcessStatus) string {
	switch status {
	case StatusRunning:
		return "🚀"
	case StatusBlocked:
		return "⏸️"
	case StatusSuspended:
		return "⏹️"
	case StatusTerminated:
		return "✅"
	default:
		return "❓"
	}
}

// GetActiveProcesses returns running processes sorted by activity
func GetActiveProcesses() ([]Process, error) {
	var processes []Process

	// Subquery to get processes with recent logs (last 7 days)
	weekAgo := time.Now().AddDate(0, 0, -7)

	err := DB.Where("status = ?", StatusRunning).
		Preload("Logs", "created_at > ?", weekAgo).
		Order("ranking DESC, updated_at DESC").
		Find(&processes).Error

	return processes, err
}

// GetBlockedProcesses returns all blocked processes with their block reasons
func GetBlockedProcesses() ([]Process, error) {
	var processes []Process

	err := DB.Where("status = ?", StatusBlocked).
		Order("created_at DESC").
		Find(&processes).Error

	return processes, err
}
