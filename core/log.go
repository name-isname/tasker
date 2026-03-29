package core

// AddLog adds a new log entry to a process
func AddLog(processID uint, logType LogType, content string) (*Log, error) {
	log := &Log{
		ProcessID: processID,
		LogType:   logType,
		Content:   content,
	}

	if err := DB.Create(log).Error; err != nil {
		return nil, err
	}

	return log, nil
}

// GetLogs retrieves all logs for a process, ordered by creation time
func GetLogs(processID uint) ([]Log, error) {
	var logs []Log
	err := DB.Where("process_id = ?", processID).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

// GetLogsPaginated retrieves logs with pagination
func GetLogsPaginated(processID uint, page, pageSize int) ([]Log, int64, error) {
	var logs []Log
	var total int64

	// Count total logs
	if err := DB.Model(&Log{}).Where("process_id = ?", processID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated logs
	offset := (page - 1) * pageSize
	err := DB.Where("process_id = ?", processID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&logs).Error

	return logs, total, err
}

// DeleteLog removes a log entry
func DeleteLog(id uint) error {
	return DB.Delete(&Log{}, id).Error
}

// UpdateLog modifies a log entry's content
func UpdateLog(id uint, content string) error {
	return DB.Model(&Log{}).Where("id = ?", id).Update("content", content).Error
}

// GetAllLogs retrieves logs across all processes with optional filtering
func GetAllLogs(logType *LogType, limit int) ([]Log, error) {
	var logs []Log
	query := DB.Order("created_at DESC")

	if logType != nil {
		query = query.Where("log_type = ?", *logType)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&logs).Error
	return logs, err
}
