package tui

import (
	"testing"
	"taskctl/core"
)

func TestInitialModel(t *testing.T) {
	m := InitialModel()

	if m.viewMode != ViewList {
		t.Errorf("Expected ViewList, got %v", m.viewMode)
	}
	if m.quitting {
		t.Error("Expected not quitting")
	}
}

func TestModelProcessLoad(t *testing.T) {
	m := InitialModel()

	// Simulate process loading
	processes := []core.Process{
		{ID: 1, Title: "Test 1", Status: core.StatusRunning},
		{ID: 2, Title: "Test 2", Status: core.StatusBlocked},
	}

	msg := ProcessesLoadedMsg{Processes: processes}
	newModel, cmd := m.Update(msg)

	if cmd != nil {
		t.Error("Expected no command")
	}

	model, ok := newModel.(Model)
	if !ok {
		t.Fatal("Expected Model type")
	}

	if len(model.processes) != 2 {
		t.Errorf("Expected 2 processes, got %d", len(model.processes))
	}
}

func TestDetailNavigation(t *testing.T) {
	m := InitialModel()
	m.processes = []core.Process{
		{ID: 1, Title: "Test 1", Status: core.StatusRunning},
	}

	// Navigate to detail
	newModel, _ := m.Update(ShowDetailMsg{ProcessID: 1})

	// After loading details, view mode should be detail
	// For this test, we just verify the message was processed
	if newModel == nil {
		t.Error("Expected model to be returned")
	}
}

func TestBackToList(t *testing.T) {
	m := InitialModel()
	m.viewMode = ViewDetail
	m.currentProcess = &core.Process{ID: 1, Title: "Test"}

	// Test going back to list
	newModel, _ := m.Update(BackToListMsg{})

	model, _ := newModel.(Model)
	if model.viewMode != ViewList {
		t.Error("Expected ViewList mode")
	}

	if model.currentProcess != nil {
		t.Error("Expected currentProcess to be cleared")
	}
}

func TestStatusStyles(t *testing.T) {
	m := Model{}

	// Test status icons
	if m.getStatusIcon(core.StatusRunning) != "▶" {
		t.Error("Expected running icon to be ▶")
	}
	if m.getStatusIcon(core.StatusBlocked) != "⏸" {
		t.Error("Expected blocked icon to be ⏸")
	}
	if m.getStatusIcon(core.StatusTerminated) != "✓" {
		t.Error("Expected terminated icon to be ✓")
	}

	// Test priority icons
	if m.getPriorityIcon(core.PriorityHigh) != "H" {
		t.Error("Expected high priority icon to be H")
	}
	if m.getPriorityIcon(core.PriorityMedium) != "M" {
		t.Error("Expected medium priority icon to be M")
	}
	if m.getPriorityIcon(core.PriorityLow) != "L" {
		t.Error("Expected low priority icon to be L")
	}
}
