package core

import (
	"time"
)

// ProcessNode represents a hierarchical process tree for visualization
type ProcessNode struct {
	Process
	Children []*ProcessNode `json:"children,omitempty"`
}

// TimelineEntry represents a log entry with its parent process context
type TimelineEntry struct {
	ID          uint      `json:"id"`
	ProcessID   uint      `json:"process_id"`
	ProcessTitle string   `json:"process_title"`
	LogType     LogType   `json:"log_type"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
}

// ActivityStat represents daily log counts for heatmaps
type ActivityStat struct {
	Date  string `json:"date"` // Format: YYYY-MM-DD
	Count int    `json:"count"`
}

// SearchResult represents a unified search result with type indicator
type SearchResult struct {
	Type     string `json:"type"` // "process" or "log"
	ID       uint   `json:"id"`
	Title    string `json:"title,omitempty"`
	Content  string `json:"content,omitempty"`
	ProcessID uint  `json:"process_id,omitempty"`
}

// ProcessExport represents complete process data for Markdown export
type ProcessExport struct {
	Process    Process
	Logs       []Log
	Children   []ProcessExport
}
