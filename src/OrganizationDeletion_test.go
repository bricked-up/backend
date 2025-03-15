package backend

import (
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

func TestDeleteOrganization(t *testing.T) {
	// Open in-memory database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	defer db.Close()

	// Execute init.sql to create the schema
	initSQL, err := os.ReadFile("../sql/init.sql")
	if err != nil {
		t.Fatalf("failed to open init.sql: %v", err)
	}
	if _, err := db.Exec(string(initSQL)); err != nil {
		t.Fatalf("failed to exec init.sql: %v", err)
	}

	// Enable foreign key support for SQLite
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
	}

	// Create a test user first
	_, err = db.Exec("INSERT INTO USER (email, password, name, verified) VALUES (?, ?, ?, ?)",
		"test@example.com", "password", "Test User", 1)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Set up a session for the user (using the proper format)
	_, err = db.Exec("INSERT INTO SESSION (userid, expires) VALUES (?, ?)", 1, "2022-01-01 00:00:00")
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	// Get the session ID
	var sessionID int
	err = db.QueryRow("SELECT id FROM SESSION WHERE userid = ?", 1).Scan(&sessionID)
	if err != nil {
		t.Fatalf("failed to get session ID: %v", err)
	}

	// Test Case 1: Create and Delete a single organization
	orgName := "Test Organization 1"
	orgID, err := CreateOrganization(db, sessionID, orgName)
	if err != nil {
		t.Fatalf("CreateOrganization returned error: %v", err)
	}

	// Now delete the organization
	err = DeleteOrganization(db, sessionID, orgID)
	if err != nil {
		t.Errorf("DeleteOrganization returned error: %v", err)
	}

}
