package core

import "errors"

// ErrCircularReference is returned when a parent-child relationship would create a cycle
var ErrCircularReference = errors.New("circular reference: parent cannot be a descendant of the child")

// wouldCreateCircularReference checks if setting parentID for processID would create a cycle
func wouldCreateCircularReference(processID, parentID uint) bool {
	// Walk up the ancestor chain from parentID to see if we reach processID
	currentID := parentID
	visited := make(map[uint]bool)

	for currentID != 0 {
		// Prevent infinite loops in case of existing corrupted data
		if visited[currentID] {
			return true // Existing cycle detected
		}
		visited[currentID] = true

		// If we found the processID in the ancestor chain, setting this parent would create a cycle
		if currentID == processID {
			return true
		}

		// Get parent of current
		var parent Process
		err := DB.Select("parent_id").First(&parent, currentID).Error
		if err != nil {
			break // Parent not found, end of chain
		}
		if parent.ParentID == nil {
			break // Reached root
		}
		currentID = *parent.ParentID
	}

	return false
}

// CreateProcess creates a new process
func CreateProcess(title, description string, parentID *uint, priority ProcessPriority) (*Process, error) {
	// Validate parentID exists if specified
	if parentID != nil {
		var parent Process
		if err := DB.First(&parent, *parentID).Error; err != nil {
			return nil, errors.New("parent process not found")
		}
	}

	process := &Process{
		Title:       title,
		Description: description,
		ParentID:    parentID,
		Status:      StatusRunning,
		Priority:    priority,
		Ranking:     0,
	}

	if err := DB.Create(process).Error; err != nil {
		return nil, err
	}

	// Auto-create initial log entry
	_, _ = AddLog(process.ID, LogTypeStateChange, "Process created")

	return process, nil
}

// GetProcess retrieves a process by ID
func GetProcess(id uint) (*Process, error) {
	var process Process
	err := DB.Preload("Parent").Preload("Logs").First(&process, id).Error
	if err != nil {
		return nil, err
	}
	return &process, nil
}

// ListProcesses returns all processes, optionally filtered by status
func ListProcesses(status *ProcessStatus) ([]Process, error) {
	var processes []Process
	query := DB.Order("ranking DESC, created_at DESC")

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	err := query.Find(&processes).Error
	return processes, err
}

// GetRootProcesses returns all top-level processes (no parent)
func GetRootProcesses() ([]Process, error) {
	var processes []Process
	err := DB.Where("parent_id IS NULL").
		Order("ranking DESC, created_at DESC").
		Find(&processes).Error
	return processes, err
}

// GetChildProcesses returns all child processes of a parent
func GetChildProcesses(parentID uint) ([]Process, error) {
	var processes []Process
	err := DB.Where("parent_id = ?", parentID).
		Order("ranking DESC, created_at DESC").
		Find(&processes).Error
	return processes, err
}

// UpdateProcess updates process fields
func UpdateProcess(id uint, title, description *string, priority *ProcessPriority, parentID *uint) error {
	updates := make(map[string]interface{})

	if title != nil {
		updates["title"] = *title
	}
	if description != nil {
		updates["description"] = *description
	}
	if priority != nil {
		updates["priority"] = *priority
	}
	if parentID != nil {
		// Check for circular reference
		if wouldCreateCircularReference(id, *parentID) {
			return ErrCircularReference
		}
		// Validate parent exists
		var parent Process
		if err := DB.First(&parent, *parentID).Error; err != nil {
			return errors.New("parent process not found")
		}
		updates["parent_id"] = *parentID
	}

	if len(updates) == 0 {
		return nil
	}

	return DB.Model(&Process{}).Where("id = ?", id).Updates(updates).Error
}

// SetProcessStatus changes the status of a process
func SetProcessStatus(id uint, status ProcessStatus) error {
	// Get current status for logging
	var process Process
	if err := DB.First(&process, id).Error; err != nil {
		return err
	}

	// Update status
	err := DB.Model(&Process{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		return err
	}

	// Log the state change
	_, _ = AddLog(id, LogTypeStateChange, "Status changed from "+string(process.Status)+" to "+string(status))

	return nil
}

// SetProcessRanking updates the ranking (sort weight) of a process
func SetProcessRanking(id uint, ranking float64) error {
	return DB.Model(&Process{}).Where("id = ?", id).Update("ranking", ranking).Error
}

// DeleteProcess removes a process and all its children recursively
func DeleteProcess(id uint) error {
	// First, get all descendant IDs recursively
	var descendantIDs []uint
	getDescendantIDs(id, &descendantIDs)

	// Delete logs for all descendants
	if len(descendantIDs) > 0 {
		DB.Where("process_id IN ?", descendantIDs).Delete(&Log{})
	}
	DB.Where("process_id = ?", id).Delete(&Log{})

	// Delete all descendants (children first due to foreign key)
	if len(descendantIDs) > 0 {
		DB.Where("id IN ?", descendantIDs).Delete(&Process{})
	}

	// Delete the process itself
	return DB.Delete(&Process{}, id).Error
}

// GetDescendantIDs returns all descendant IDs of a process (for filtering)
func GetDescendantIDs(parentID uint) ([]uint, error) {
	var ids []uint
	getDescendantIDs(parentID, &ids)
	return ids, nil
}

// getDescendantIDs recursively collects all descendant IDs
func getDescendantIDs(parentID uint, ids *[]uint) {
	getDescendantIDsHelper(parentID, ids, make(map[uint]bool))
}

// getDescendantIDsHelper recursively collects descendant IDs with cycle detection
func getDescendantIDsHelper(parentID uint, ids *[]uint, visited map[uint]bool) {
	// Prevent infinite recursion due to circular references
	if visited[parentID] {
		return
	}
	visited[parentID] = true

	var children []Process
	DB.Where("parent_id = ?", parentID).Pluck("id", &children)

	for _, child := range children {
		*ids = append(*ids, child.ID)
		getDescendantIDsHelper(child.ID, ids, visited) // Recurse
	}
}
