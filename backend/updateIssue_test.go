package backend

import (
	"testing"

	_ "modernc.org/sqlite"
)

// TestUpdateIssue verifies the UpdateIssue function for various conditions
func TestUpdateIssue(t *testing.T) {
	db := setupDatabase(t)
	defer db.Close()

	// Test case for successful issue update
	t.Run("Successful issue update", func(t *testing.T) {
		completed := "2025-04-01"
		err := UpdateIssue(db, 1, 1, "Updated Title", "Updated Description", "2025-01-01", &completed, 1000)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	// Test case for invalid session ID
	t.Run("Invalid session ID", func(t *testing.T) {
		completed := "2025-04-01"
		err := UpdateIssue(db, 999, 1, "Updated Title", "Updated Description", "2025-01-01", &completed, 1000)
		if err == nil || err.Error() != "Invalid session ID" {
			t.Errorf("Expected error 'Invalid session ID', got %v", err)
		}
	})

	// Test case for invalid issue ID
	t.Run("Invalid issue ID", func(t *testing.T) {
		completed := "2025-04-01"
		err := UpdateIssue(db, 1, 999, "Updated Title", "Updated Description", "2025-01-01", &completed, 1000)
		if err == nil || err.Error() != "Invalid issue ID" {
			t.Errorf("Expected error 'Invalid issue ID', got %v", err)
		}
	})

	// Test case for user without write permissions
	t.Run("User without write permissions", func(t *testing.T) {
		_, err := db.Exec("UPDATE PROJECT_ROLE SET can_write = 0 WHERE id = 1")
		if err != nil {
			t.Fatalf("Failed to modify permissions: %v", err)
		}

		completed := "2025-04-01"
		err = UpdateIssue(db, 4, 1, "Updated Title", "Updated Description", "2025-01-01", &completed, 1000)
		if err == nil || err.Error() != "User does not have write permissions" {
			t.Errorf("Expected error 'User does not have write permissions', got %v", err)
		}
	})
}
