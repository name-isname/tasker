package core

import (
	"os"
	"path/filepath"
	"gorm.io/gorm"
	"github.com/glebarez/sqlite"
)

var DB *gorm.DB

// InitDB initializes the SQLite database
func InitDB(dbPath string) error {
	// Ensure parent directory exists
	dbDir := filepath.Dir(dbPath)
	if dbDir != "." && dbDir != "" {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return err
		}
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

// AutoMigrate runs database migrations and creates FTS tables
func AutoMigrate() error {
	// Migrate main schemas
	if err := DB.AutoMigrate(&Process{}, &Log{}); err != nil {
		return err
	}

	// Create FTS5 virtual table for full-text search
	// Using pure SQL since GORM doesn't support virtual tables directly
	if err := createFTSTable(); err != nil {
		return err
	}

	// Create triggers to keep FTS table in sync
	if err := createFTSTriggers(); err != nil {
		return err
	}

	return nil
}

// createFTSTable creates the FTS5 virtual table
func createFTSTable() error {
	// Drop existing table if any (for idempotency during development)
	DB.Exec("DROP TABLE IF EXISTS process_fts")

	// Create FTS5 virtual table
	// We store id, title, and a combined content field for searching
	sql := `
		CREATE VIRTUAL TABLE process_fts USING fts5(
			id UNINDEXED,
			title,
			content
		);
	`
	return DB.Exec(sql).Error
}

// createFTSTriggers creates triggers to keep FTS table synchronized
func createFTSTriggers() error {
	// Drop existing triggers
	DB.Exec("DROP TRIGGER IF EXISTS process_fts_insert")
	DB.Exec("DROP TRIGGER IF EXISTS process_fts_update")
	DB.Exec("DROP TRIGGER IF EXISTS process_fts_delete")

	// Trigger: Insert new process into FTS
	insertTrigger := `
		CREATE TRIGGER process_fts_insert AFTER INSERT ON processes BEGIN
			INSERT INTO process_fts(id, title, content)
			VALUES (NEW.id, NEW.title, NEW.title || ' ' || COALESCE(NEW.description, ''));
		END;
	`

	// Trigger: Update FTS on process change
	updateTrigger := `
		CREATE TRIGGER process_fts_update AFTER UPDATE OF title, description ON processes BEGIN
			UPDATE process_fts
			SET title = NEW.title,
			    content = NEW.title || ' ' || COALESCE(NEW.description, '')
			WHERE id = NEW.id;
		END;
	`

	// Trigger: Delete from FTS on process deletion
	deleteTrigger := `
		CREATE TRIGGER process_fts_delete AFTER DELETE ON processes BEGIN
			DELETE FROM process_fts WHERE id = OLD.id;
		END;
	`

	if err := DB.Exec(insertTrigger).Error; err != nil {
		return err
	}
	if err := DB.Exec(updateTrigger).Error; err != nil {
		return err
	}
	return DB.Exec(deleteTrigger).Error
}

// SearchProcesses performs full-text search across process titles and descriptions
func SearchProcesses(query string) ([]Process, error) {
	if query == "" {
		return []Process{}, nil
	}

	var results []ProcessFTS
	// Use FTS5 simple query syntax
	err := DB.Raw(`
		SELECT id, title, content
		FROM process_fts
		WHERE process_fts MATCH ?
		ORDER BY rank
	`, query).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Fetch full Process objects for the matching IDs
	var processIDs []uint
	for _, r := range results {
		processIDs = append(processIDs, r.ID)
	}

	if len(processIDs) == 0 {
		return []Process{}, nil
	}

	var processes []Process
	err = DB.Where("id IN ?", processIDs).Find(&processes).Error
	return processes, err
}
