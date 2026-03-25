package core

import (
	"testing"
)

func TestAddLog(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	process, _ := CreateProcess("Test", "", nil, PriorityMedium)

	log, err := AddLog(process.ID, LogTypeProgress, "Making progress")
	if err != nil {
		t.Fatalf("AddLog failed: %v", err)
	}

	if log.ID == 0 {
		t.Error("Expected non-zero log ID")
	}
	if log.ProcessID != process.ID {
		t.Errorf("Expected process ID %d, got %d", process.ID, log.ProcessID)
	}
	if log.Content != "Making progress" {
		t.Errorf("Expected content 'Making progress', got '%s'", log.Content)
	}
	if log.LogType != LogTypeProgress {
		t.Errorf("Expected log type '%s', got '%s'", LogTypeProgress, log.LogType)
	}
}

func TestGetLogs(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	process, _ := CreateProcess("Test", "", nil, PriorityMedium)

	AddLog(process.ID, LogTypeProgress, "First log")
	AddLog(process.ID, LogTypeStateChange, "Second log")
	AddLog(process.ID, LogTypeProgress, "Third log")

	logs, err := GetLogs(process.ID)
	if err != nil {
		t.Fatalf("GetLogs failed: %v", err)
	}
	if len(logs) != 4 { // 3 added + 1 auto-created on process creation
		t.Errorf("Expected 4 logs, got %d", len(logs))
	}

	// Check ordering (should be DESC by created_at)
	if logs[0].Content != "Third log" {
		t.Errorf("Expected first log to be 'Third log', got '%s'", logs[0].Content)
	}
}

func TestGetLogsPaginated(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	process, _ := CreateProcess("Test", "", nil, PriorityMedium)

	for i := 0; i < 15; i++ {
		AddLog(process.ID, LogTypeProgress, "Log entry")
	}

	logs, total, err := GetLogsPaginated(process.ID, 1, 10)
	if err != nil {
		t.Fatalf("GetLogsPaginated failed: %v", err)
	}
	if total != 16 { // 15 + 1 auto-created
		t.Errorf("Expected total 16, got %d", total)
	}
	if len(logs) != 10 {
		t.Errorf("Expected 10 logs on page 1, got %d", len(logs))
	}

	// Get second page
	logs2, total2, err := GetLogsPaginated(process.ID, 2, 10)
	if err != nil {
		t.Fatalf("GetLogsPaginated page 2 failed: %v", err)
	}
	if total2 != 16 {
		t.Errorf("Expected total 16 on page 2, got %d", total2)
	}
	if len(logs2) != 6 {
		t.Errorf("Expected 6 logs on page 2, got %d", len(logs2))
	}
}

func TestDeleteLog(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	process, _ := CreateProcess("Test", "", nil, PriorityMedium)
	log, _ := AddLog(process.ID, LogTypeProgress, "To be deleted")

	err := DeleteLog(log.ID)
	if err != nil {
		t.Fatalf("DeleteLog failed: %v", err)
	}

	logs, _ := GetLogs(process.ID)
	for _, l := range logs {
		if l.ID == log.ID {
			t.Error("Expected log to be deleted")
		}
	}
}

func TestUpdateLog(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	process, _ := CreateProcess("Test", "", nil, PriorityMedium)
	log, _ := AddLog(process.ID, LogTypeProgress, "Original content")

	newContent := "Updated content"
	err := UpdateLog(log.ID, newContent)
	if err != nil {
		t.Fatalf("UpdateLog failed: %v", err)
	}

	// Re-fetch to verify
	logs, _ := GetLogs(process.ID)
	var updated *Log
	for _, l := range logs {
		if l.ID == log.ID {
			updated = &l
			break
		}
	}

	if updated == nil {
		t.Fatal("Expected to find updated log")
	}
	if updated.Content != newContent {
		t.Errorf("Expected content '%s', got '%s'", newContent, updated.Content)
	}
}

func TestGetAllLogs(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	p1, _ := CreateProcess("P1", "", nil, PriorityMedium)
	p2, _ := CreateProcess("P2", "", nil, PriorityMedium)

	AddLog(p1.ID, LogTypeProgress, "P1 progress")
	AddLog(p1.ID, LogTypeStateChange, "P1 state")
	AddLog(p2.ID, LogTypeProgress, "P2 progress")

	logs, err := GetAllLogs(nil, 0)
	if err != nil {
		t.Fatalf("GetAllLogs failed: %v", err)
	}
	if len(logs) != 5 { // 3 added + 2 auto-created
		t.Errorf("Expected 5 logs, got %d", len(logs))
	}

	// Filter by type
	progressType := LogTypeProgress
	progressLogs, err := GetAllLogs(&progressType, 0)
	if err != nil {
		t.Fatalf("GetAllLogs with filter failed: %v", err)
	}
	if len(progressLogs) != 2 {
		t.Errorf("Expected 2 progress logs, got %d", len(progressLogs))
	}

	// Test limit
	limitedLogs, err := GetAllLogs(nil, 2)
	if err != nil {
		t.Fatalf("GetAllLogs with limit failed: %v", err)
	}
	if len(limitedLogs) != 2 {
		t.Errorf("Expected 2 logs with limit, got %d", len(limitedLogs))
	}
}

func TestProcessWithLogs(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	process, _ := CreateProcess("Test Process", "Description", nil, PriorityMedium)
	AddLog(process.ID, LogTypeProgress, "Step 1 completed")
	AddLog(process.ID, LogTypeProgress, "Step 2 completed")

	// Get process with preloaded logs
	fetched, err := GetProcess(process.ID)
	if err != nil {
		t.Fatalf("GetProcess failed: %v", err)
	}

	if len(fetched.Logs) != 3 { // 2 added + 1 auto-created
		t.Errorf("Expected 3 logs, got %d", len(fetched.Logs))
	}
}
