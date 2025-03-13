package backend

import (
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

// setupTestDB initializes an in-memory SQLite database for testing.
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		t.Fatalf("Failed to enable foreign keys: %v", err)
	}

	// Execute init.sql
	initSQL, err := os.ReadFile("../sql/init.sql")
	if err != nil {
		t.Fatalf("Failed to open init.sql: %v", err)
	}
	if _, err := db.Exec(string(initSQL)); err != nil {
		t.Fatalf("Failed to execute init.sql: %v", err)
	}

	// Execute populate.sql
	populateSQL, err := os.ReadFile("../sql/populate.sql")
	if err != nil {
		t.Fatalf("Failed to open populate.sql: %v", err)
	}
	if _, err := db.Exec(string(populateSQL)); err != nil {
		t.Fatalf("Failed to execute populate.sql: %v", err)
	}

	return db // Do not close the DB here; let the test handle closing.
}

// TestDeleteUser verifies that a user and related records are deleted correctly.
func TestDeleteUser(t *testing.T) {
	db := setupTestDB(t)
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
	err = deleteUser(db, sessionID)
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
