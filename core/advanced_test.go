package core

import (
	"testing"
	"time"
)

func TestChangeProcessState(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	process, _ := CreateProcess("Test", "", nil, PriorityMedium)

	err := ChangeProcessState(process.ID, StatusBlocked, "Waiting for approval")
	if err != nil {
		t.Fatalf("ChangeProcessState failed: %v", err)
	}

	// Verify status changed
	updated, _ := GetProcess(process.ID)
	if updated.Status != StatusBlocked {
		t.Errorf("Expected status %s, got %s", StatusBlocked, updated.Status)
	}

	// Verify log was created with reason
	logs, _ := GetLogs(process.ID)
	found := false
	for _, log := range logs {
		if log.LogType == LogTypeStateChange && log.Content == "Status changed from running to blocked. Reason: Waiting for approval" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected state change log with reason not found")
	}
}

func TestGetProcessTree(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	root, _ := CreateProcess("Root", "", nil, PriorityMedium)
	child1, _ := CreateProcess("Child1", "", &root.ID, PriorityMedium)
	_, _ = CreateProcess("Child2", "", &root.ID, PriorityMedium)
	CreateProcess("Grandchild", "", &child1.ID, PriorityMedium)

	tree, err := GetProcessTree(root.ID)
	if err != nil {
		t.Fatalf("GetProcessTree failed: %v", err)
	}

	if tree.ID != root.ID {
		t.Errorf("Expected root ID %d, got %d", root.ID, tree.ID)
	}
	if len(tree.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(tree.Children))
	}
	// Find child1 and check its grandchild
	foundGrandchild := false
	for _, child := range tree.Children {
		if len(child.Children) > 0 {
			foundGrandchild = true
			break
		}
	}
	if !foundGrandchild {
		t.Error("Expected to find grandchild in one of the children")
	}
}

func TestGetFullProcessTree(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	CreateProcess("Root1", "", nil, PriorityMedium)
	CreateProcess("Root2", "", nil, PriorityMedium)

	trees, err := GetFullProcessTree()
	if err != nil {
		t.Fatalf("GetFullProcessTree failed: %v", err)
	}
	if len(trees) != 2 {
		t.Errorf("Expected 2 root trees, got %d", len(trees))
	}
}

func TestGetTimeline(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	p1, _ := CreateProcess("P1", "", nil, PriorityMedium)
	p2, _ := CreateProcess("P2", "", nil, PriorityMedium)

	AddLog(p1.ID, LogTypeProgress, "Log 1")
	time.Sleep(time.Millisecond)
	AddLog(p2.ID, LogTypeProgress, "Log 2")
	AddLog(p1.ID, LogTypeProgress, "Log 3")

	entries, err := GetTimeline(time.Time{}, time.Time{}, 0)
	if err != nil {
		t.Fatalf("GetTimeline failed: %v", err)
	}

	// Should have 3 progress logs + 2 state change logs from creation
	if len(entries) < 3 {
		t.Errorf("Expected at least 3 timeline entries, got %d", len(entries))
	}

	// Check ordering (should be DESC)
	if entries[0].Content != "Log 3" {
		t.Errorf("Expected first entry to be 'Log 3', got '%s'", entries[0].Content)
	}
}

func TestGetTodayTimeline(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	process, _ := CreateProcess("Test", "", nil, PriorityMedium)
	AddLog(process.ID, LogTypeProgress, "Today's log")

	entries, err := GetTodayTimeline()
	if err != nil {
		t.Fatalf("GetTodayTimeline failed: %v", err)
	}

	if len(entries) == 0 {
		t.Error("Expected at least one entry for today")
	}
}

func TestGetActivityStats(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	p1, _ := CreateProcess("P1", "", nil, PriorityMedium)
	p2, _ := CreateProcess("P2", "", nil, PriorityMedium)

	AddLog(p1.ID, LogTypeProgress, "Log 1")
	AddLog(p1.ID, LogTypeProgress, "Log 2")
	AddLog(p2.ID, LogTypeProgress, "Log 3")

	stats, err := GetActivityStats(7)
	if err != nil {
		t.Fatalf("GetActivityStats failed: %v", err)
	}

	if len(stats) == 0 {
		t.Error("Expected activity stats")
	}

	// Check today's count
	today := time.Now().Format("2006-01-02")
	foundToday := false
	for _, stat := range stats {
		if stat.Date == today {
			foundToday = true
			// Should have at least 3 progress logs + 2 state change logs
			if stat.Count < 3 {
				t.Errorf("Expected at least 3 logs today, got %d", stat.Count)
			}
		}
	}
	if !foundToday {
		t.Error("Expected to find today's activity")
	}
}

func TestGlobalSearch(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	CreateProcess("golang programming", "Learn Go", nil, PriorityMedium)
	CreateProcess("python basics", "Learn Python", nil, PriorityMedium)

	p1, _ := CreateProcess("Test", "", nil, PriorityMedium)
	AddLog(p1.ID, LogTypeProgress, "Fixed database bug")

	results, err := GlobalSearch("program")
	if err != nil {
		t.Fatalf("GlobalSearch failed: %v", err)
	}

	if len(results) == 0 {
		t.Error("Expected search results for 'program'")
	}

	// Test log search
	results, err = GlobalSearch("database")
	if err != nil {
		t.Fatalf("GlobalSearch log failed: %v", err)
	}

	foundLog := false
	for _, r := range results {
		if r.Type == "log" {
			foundLog = true
			break
		}
	}
	if !foundLog {
		t.Error("Expected to find log containing 'database'")
	}
}

func TestGetProcessContext(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	parent, _ := CreateProcess("Parent", "Parent desc", nil, PriorityMedium)
	child, _ := CreateProcess("Child", "Child desc", &parent.ID, PriorityMedium)

	AddLog(parent.ID, LogTypeProgress, "Parent log 1")
	AddLog(parent.ID, LogTypeProgress, "Parent log 2")
	AddLog(child.ID, LogTypeProgress, "Child log")

	ctx, err := GetProcessContext(parent.ID)
	if err != nil {
		t.Fatalf("GetProcessContext failed: %v", err)
	}

	if ctx.Process.ID != parent.ID {
		t.Errorf("Expected process ID %d, got %d", parent.ID, ctx.Process.ID)
	}

	// Should have 2 parent logs + 1 state change log
	if len(ctx.Logs) != 3 {
		t.Errorf("Expected 3 logs, got %d", len(ctx.Logs))
	}

	// Should have 1 child
	if len(ctx.Children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(ctx.Children))
	}
}

func TestGetActiveProcesses(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	p1, _ := CreateProcess("Active 1", "", nil, PriorityMedium)
	p2, _ := CreateProcess("Active 2", "", nil, PriorityMedium)

	// Make one blocked
	SetProcessStatus(p2.ID, StatusBlocked)

	processes, err := GetActiveProcesses()
	if err != nil {
		t.Fatalf("GetActiveProcesses failed: %v", err)
	}

	if len(processes) != 1 {
		t.Errorf("Expected 1 active process, got %d", len(processes))
	}
	if processes[0].ID != p1.ID {
		t.Errorf("Expected process %d, got %d", p1.ID, processes[0].ID)
	}
}

func TestGetBlockedProcesses(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	_, _ = CreateProcess("Running", "", nil, PriorityMedium)
	p2, _ := CreateProcess("Blocked", "", nil, PriorityMedium)

	SetProcessStatus(p2.ID, StatusBlocked)

	processes, err := GetBlockedProcesses()
	if err != nil {
		t.Fatalf("GetBlockedProcesses failed: %v", err)
	}

	if len(processes) != 1 {
		t.Errorf("Expected 1 blocked process, got %d", len(processes))
	}
	if processes[0].ID != p2.ID {
		t.Errorf("Expected process %d, got %d", p2.ID, processes[0].ID)
	}
}
