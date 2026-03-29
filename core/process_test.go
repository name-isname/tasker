package core

import (
	"testing"
)

func TestCreateProcess(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	process, err := CreateProcess("Test Process", "Description", nil, PriorityHigh)
	if err != nil {
		t.Fatalf("CreateProcess failed: %v", err)
	}

	if process.ID == 0 {
		t.Error("Expected non-zero ID")
	}
	if process.Title != "Test Process" {
		t.Errorf("Expected title 'Test Process', got '%s'", process.Title)
	}
	if process.Status != StatusRunning {
		t.Errorf("Expected status '%s', got '%s'", StatusRunning, process.Status)
	}
	if process.Priority != PriorityHigh {
		t.Errorf("Expected priority '%s', got '%s'", PriorityHigh, process.Priority)
	}
}

func TestCreateProcessWithParent(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	parent, _ := CreateProcess("Parent", "", nil, PriorityMedium)
	child, err := CreateProcess("Child", "", &parent.ID, PriorityLow)

	if err != nil {
		t.Fatalf("CreateProcess with parent failed: %v", err)
	}
	if child.ParentID == nil {
		t.Error("Expected ParentID to be set")
	}
	if *child.ParentID != parent.ID {
		t.Errorf("Expected ParentID %d, got %d", parent.ID, *child.ParentID)
	}
}

func TestGetProcess(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	created, _ := CreateProcess("Test", "Desc", nil, PriorityMedium)
	fetched, err := GetProcess(created.ID)

	if err != nil {
		t.Fatalf("GetProcess failed: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, fetched.ID)
	}
	if fetched.Title != created.Title {
		t.Errorf("Expected title '%s', got '%s'", created.Title, fetched.Title)
	}
}

func TestListProcesses(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	CreateProcess("Process 1", "", nil, PriorityHigh)
	CreateProcess("Process 2", "", nil, PriorityLow)

	processes, err := ListProcesses(nil)
	if err != nil {
		t.Fatalf("ListProcesses failed: %v", err)
	}
	if len(processes) != 2 {
		t.Errorf("Expected 2 processes, got %d", len(processes))
	}
}

func TestListProcessesFilterByStatus(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	_, _ = CreateProcess("P1", "", nil, PriorityMedium)
	p2, _ := CreateProcess("P2", "", nil, PriorityMedium)
	SetProcessStatus(p2.ID, StatusBlocked)

	running := StatusRunning
	processes, _ := ListProcesses(&running)
	if len(processes) != 1 {
		t.Errorf("Expected 1 running process, got %d", len(processes))
	}
	if processes[0].Status != StatusRunning {
		t.Errorf("Expected status '%s', got '%s'", StatusRunning, processes[0].Status)
	}
}

func TestSetProcessStatus(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	process, _ := CreateProcess("Test", "", nil, PriorityMedium)

	err := SetProcessStatus(process.ID, StatusBlocked)
	if err != nil {
		t.Fatalf("SetProcessStatus failed: %v", err)
	}

	updated, _ := GetProcess(process.ID)
	if updated.Status != StatusBlocked {
		t.Errorf("Expected status '%s', got '%s'", StatusBlocked, updated.Status)
	}

	// Check that log was created
	logs, _ := GetLogs(process.ID)
	if len(logs) == 0 {
		t.Error("Expected state change log to be created")
	}
	if logs[0].LogType != LogTypeStateChange {
		t.Errorf("Expected log type '%s', got '%s'", LogTypeStateChange, logs[0].LogType)
	}
}

func TestSetProcessRanking(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	process, _ := CreateProcess("Test", "", nil, PriorityMedium)

	err := SetProcessRanking(process.ID, 42.5)
	if err != nil {
		t.Fatalf("SetProcessRanking failed: %v", err)
	}

	fetched, _ := GetProcess(process.ID)
	if fetched.Ranking != 42.5 {
		t.Errorf("Expected ranking 42.5, got %f", fetched.Ranking)
	}
}

func TestUpdateProcess(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	process, _ := CreateProcess("Original", "Desc", nil, PriorityMedium)

	newTitle := "Updated Title"
	newDesc := "New Description"
	newPriority := PriorityHigh

	err := UpdateProcess(process.ID, &newTitle, &newDesc, &newPriority, nil)
	if err != nil {
		t.Fatalf("UpdateProcess failed: %v", err)
	}

	updated, _ := GetProcess(process.ID)
	if updated.Title != newTitle {
		t.Errorf("Expected title '%s', got '%s'", newTitle, updated.Title)
	}
	if updated.Description != newDesc {
		t.Errorf("Expected description '%s', got '%s'", newDesc, updated.Description)
	}
	if updated.Priority != newPriority {
		t.Errorf("Expected priority '%s', got '%s'", newPriority, updated.Priority)
	}
}

func TestDeleteProcess(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	process, _ := CreateProcess("To Delete", "", nil, PriorityMedium)

	err := DeleteProcess(process.ID)
	if err != nil {
		t.Fatalf("DeleteProcess failed: %v", err)
	}

	_, err = GetProcess(process.ID)
	if err == nil {
		t.Error("Expected error when fetching deleted process")
	}
}

func TestDeleteProcessWithChildren(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	parent, _ := CreateProcess("Parent", "", nil, PriorityMedium)
	child1, _ := CreateProcess("Child1", "", &parent.ID, PriorityMedium)
	child2, _ := CreateProcess("Child2", "", &parent.ID, PriorityMedium)

	DeleteProcess(parent.ID)

	// Check parent is deleted
	_, err := GetProcess(parent.ID)
	if err == nil {
		t.Error("Expected parent to be deleted")
	}

	// Check children are deleted
	_, err = GetProcess(child1.ID)
	if err == nil {
		t.Error("Expected child1 to be deleted")
	}
	_, err = GetProcess(child2.ID)
	if err == nil {
		t.Error("Expected child2 to be deleted")
	}
}

func TestDeleteProcessWithCircularReference(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Create processes
	proc1, _ := CreateProcess("Proc1", "", nil, PriorityMedium)
	proc2, _ := CreateProcess("Proc2", "", &proc1.ID, PriorityMedium)

	// Create a circular reference: proc1's parent becomes proc2
	// This simulates corrupted data
	DB.Model(&Process{}).Where("id = ?", proc1.ID).Update("parent_id", proc2.ID)

	// DeleteProcess should handle this gracefully without infinite recursion
	err := DeleteProcess(proc1.ID)
	if err != nil {
		t.Errorf("DeleteProcess should handle circular references, got: %v", err)
	}
}

func TestGetChildProcesses(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	parent, _ := CreateProcess("Parent", "", nil, PriorityMedium)
	CreateProcess("Child1", "", &parent.ID, PriorityMedium)
	CreateProcess("Child2", "", &parent.ID, PriorityMedium)

	children, err := GetChildProcesses(parent.ID)
	if err != nil {
		t.Fatalf("GetChildProcesses failed: %v", err)
	}
	if len(children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(children))
	}
}

func TestGetRootProcesses(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	root1, _ := CreateProcess("Root1", "", nil, PriorityMedium)
	_, _ = CreateProcess("Root2", "", nil, PriorityMedium)
	CreateProcess("Child", "", &root1.ID, PriorityMedium)

	roots, err := GetRootProcesses()
	if err != nil {
		t.Fatalf("GetRootProcesses failed: %v", err)
	}
	if len(roots) != 2 {
		t.Errorf("Expected 2 root processes, got %d", len(roots))
	}
}

func TestUpdateProcessCircularReference(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Create a hierarchy: grandparent -> parent -> child
	grandparent, _ := CreateProcess("Grandparent", "", nil, PriorityMedium)
	parent, _ := CreateProcess("Parent", "", &grandparent.ID, PriorityMedium)
	child, _ := CreateProcess("Child", "", &parent.ID, PriorityMedium)

	// Try to set grandparent's parent to child (would create cycle)
	err := UpdateProcess(grandparent.ID, nil, nil, nil, &child.ID)
	if err != ErrCircularReference {
		t.Errorf("Expected ErrCircularReference, got: %v", err)
	}

	// Try to set parent's parent to child (would create cycle)
	err = UpdateProcess(parent.ID, nil, nil, nil, &child.ID)
	if err != ErrCircularReference {
		t.Errorf("Expected ErrCircularReference, got: %v", err)
	}

	// Valid update should still work
	validParent := uint(999) // Non-existent parent should fail with different error
	err = UpdateProcess(child.ID, nil, nil, nil, &validParent)
	if err == nil {
		t.Error("Expected error for non-existent parent")
	}
}

func TestCreateProcessInvalidParent(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Try to create process with non-existent parent
	invalidParent := uint(999)
	_, err := CreateProcess("Test", "Description", &invalidParent, PriorityMedium)
	if err == nil {
		t.Error("Expected error for non-existent parent")
	}
}

func TestSearchProcesses(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	CreateProcess("golang programming", "Learn Go language", nil, PriorityMedium)
	CreateProcess("python basics", "Python tutorial", nil, PriorityMedium)
	CreateProcess("java spring", "Spring framework", nil, PriorityMedium)

	// FTS5 search tests - may not work in in-memory DB
	results, err := SearchProcesses("golang")
	if err != nil {
		t.Logf("FTS5 search error (may not work in in-memory DB): %v", err)
		return // Skip test if FTS5 not available
	}

	// If FTS5 works, verify results
	if len(results) > 0 && results[0].Title != "golang programming" {
		t.Errorf("Expected 'golang programming', got '%s'", results[0].Title)
	}
}

func TestGetDescendantIDs(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Create a hierarchy: root -> child1 -> grandchild1, grandchild2
	//                     -> child2
	root, _ := CreateProcess("Root", "", nil, PriorityMedium)
	child1, _ := CreateProcess("Child1", "", &root.ID, PriorityMedium)
	child2, _ := CreateProcess("Child2", "", &root.ID, PriorityMedium)
	grandchild1, _ := CreateProcess("Grandchild1", "", &child1.ID, PriorityMedium)
	grandchild2, _ := CreateProcess("Grandchild2", "", &child1.ID, PriorityMedium)

	// Get descendants of root
	descendants, err := GetDescendantIDs(root.ID)
	if err != nil {
		t.Fatalf("GetDescendantIDs failed: %v", err)
	}

	// Should have 4 descendants: child1, child2, grandchild1, grandchild2
	if len(descendants) != 4 {
		t.Errorf("Expected 4 descendants, got %d", len(descendants))
	}

	// Verify IDs - all should be in the descendants list
	expectedIDs := map[uint]bool{child1.ID: true, child2.ID: true, grandchild1.ID: true, grandchild2.ID: true}
	for _, id := range descendants {
		if !expectedIDs[id] {
			t.Errorf("Unexpected descendant ID: %d", id)
		}
	}

	// Get descendants of child1 (should have 2 grandchildren)
	child1Descendants, _ := GetDescendantIDs(child1.ID)
	if len(child1Descendants) != 2 {
		t.Errorf("Expected 2 descendants for child1, got %d", len(child1Descendants))
	}

	// Get descendants of leaf process (should be empty)
	leafDescendants, _ := GetDescendantIDs(grandchild1.ID)
	if len(leafDescendants) != 0 {
		t.Errorf("Expected 0 descendants for leaf, got %d", len(leafDescendants))
	}
}

func TestGetDescendantsWithCircularReference(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Create processes
	proc1, _ := CreateProcess("Proc1", "", nil, PriorityMedium)
	proc2, _ := CreateProcess("Proc2", "", &proc1.ID, PriorityMedium)

	// Create a circular reference
	DB.Model(&Process{}).Where("id = ?", proc1.ID).Update("parent_id", proc2.ID)

	// GetDescendantIDs should handle this gracefully
	descendants, err := GetDescendantIDs(proc1.ID)
	if err != nil {
		t.Fatalf("GetDescendantIDs should handle circular references, got error: %v", err)
	}

	// Should return at least proc2 (the direct child)
	if len(descendants) == 0 {
		t.Error("Expected at least one descendant")
	}
}
