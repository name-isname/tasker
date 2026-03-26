package core

// CreateProcess creates a new process
func CreateProcess(title, description string, parentID *uint, priority ProcessPriority) (*Process, error) {
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

// getDescendantIDs recursively collects all descendant IDs
func getDescendantIDs(parentID uint, ids *[]uint) {
	var children []Process
	DB.Where("parent_id = ?", parentID).Pluck("id", &children)

	for _, child := range children {
		*ids = append(*ids, child.ID)
		getDescendantIDs(child.ID, ids) // Recurse
	}
}
