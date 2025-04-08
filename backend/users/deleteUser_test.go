package users

import (
	"brickedup/backend/utils"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

// TestDeleteUser verifies that a user and related records are deleted correctly.
func TestDeleteUser(t *testing.T) {
	db := utils.SetupTest(t)
	defer db.Close()

	// Fetch a session ID for an existing user
	var sessionID string
	err := db.QueryRow("SELECT id FROM SESSION LIMIT 1").Scan(&sessionID)
	if err == sql.ErrNoRows {
		t.Fatalf("No session found in SESSION table")
	} else if err != nil {
		t.Fatalf("Failed to fetch session ID: %v", err)
	}

	// Attempt to delete the user
	err = DeleteUser(db, sessionID)
	if err != nil {
		t.Fatalf("deleteUser returned an error: %v", err)
	}

	// Verify that the user is deleted
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM USER WHERE id = (SELECT userid FROM SESSION WHERE id = ?)", sessionID).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query user count: %v", err)
	}

	if count != 0 {
		t.Errorf("User was not deleted, count: %d", count)
	}
}
