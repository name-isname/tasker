package core

import (
	"time"
)

// ProcessStatus represents the current state of a process
type ProcessStatus string

const (
	StatusRunning    ProcessStatus = "running"
	StatusBlocked    ProcessStatus = "blocked"
	StatusSuspended  ProcessStatus = "suspended"
	StatusTerminated ProcessStatus = "terminated"
)

// ProcessPriority represents the priority level
type ProcessPriority string

const (
	PriorityLow    ProcessPriority = "low"
	PriorityMedium ProcessPriority = "medium"
	PriorityHigh   ProcessPriority = "high"
)

// Process represents a task/goal in the system
// It models tasks as OS processes with state transitions
type Process struct {
	ID          uint             `json:"id" gorm:"primaryKey"`
	ParentID    *uint            `json:"parent_id,omitempty" gorm:"index"`
	Title       string           `json:"title" gorm:"not null;size:255"`
	Description string           `json:"description" gorm:"type:text"`
	Status      ProcessStatus    `json:"status" gorm:"not null;size:20;default:running"`
	Priority    ProcessPriority  `json:"priority" gorm:"not null;size:10;default:medium"`
	Ranking     float64          `json:"ranking" gorm:"default:0"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`

	// Relationships
	Parent   *Process   `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Process  `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Logs     []Log      `json:"logs,omitempty" gorm:"foreignKey:ProcessID"`
}

// LogType represents the type of log entry
type LogType string

const (
	LogTypeStateChange LogType = "state_change"
	LogTypeProgress    LogType = "progress"
)

// Log represents a timeline entry for a process
// It records incremental progress, roadblocks, and state changes
type Log struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProcessID uint      `json:"process_id" gorm:"not null;index"`
	LogType   LogType   `json:"log_type" gorm:"not null;size:20"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at"`

	// Relationship
	Process Process `json:"-" gorm:"foreignKey:ProcessID"`
}

// ProcessFTS represents the full-text search table for processes
// This uses SQLite FTS5 for fast global search
type ProcessFTS struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Title string `gorm:"type:text" json:"title"`
	// Note: Content combines title, description, and logs for FTS
	Content string `gorm:"type:text" json:"content"`
}

// TableName specifies the FTS virtual table name
func (ProcessFTS) TableName() string {
	return "process_fts"
}
