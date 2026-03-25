package core

import (
	"testing"
)

// setupTestDB initializes an in-memory database for testing
func setupTestDB(t *testing.T) {
	err := InitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to init test DB: %v", err)
	}
	err = AutoMigrate()
	if err != nil {
		t.Fatalf("Failed to migrate test DB: %v", err)
	}
}

// teardownTestDB cleans up the test database
func teardownTestDB() {
	DB = nil
}
