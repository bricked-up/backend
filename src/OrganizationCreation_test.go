package backend

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestCreateOrganization(t *testing.T) {
	// Open in-memory database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	defer db.Close()

	// Create necessary tables
	_, err = db.Exec(`
		CREATE TABLE sessions (session_id TEXT PRIMARY KEY, user_id INTEGER);
		CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT);
		CREATE TABLE organizations (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT UNIQUE);
		CREATE TABLE organization_roles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			organization_id INTEGER,
			role_name TEXT,
			can_read INTEGER,
			can_write INTEGER,
			can_execute INTEGER
		);
		CREATE TABLE organization_members (user_id INTEGER, organization_id INTEGER, role_id INTEGER);
	`)
	if err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	// Insert a test user and session
	_, err = db.Exec("INSERT INTO users (username) VALUES ('testuser')")
	if err != nil {
		t.Fatalf("failed to insert test user: %v", err)
	}

	var userID int
	err = db.QueryRow("SELECT id FROM users WHERE username = 'testuser'").Scan(&userID)
	if err != nil {
		t.Fatalf("failed to get test user ID: %v", err)
	}

	sessionID := "testsession123"
	_, err = db.Exec("INSERT INTO sessions (session_id, user_id) VALUES (?, ?)", sessionID, userID)
	if err != nil {
		t.Fatalf("failed to insert test session: %v", err)
	}

	// Test valid organization creation
	orgName := "TestOrg"
	orgID, err := CreateOrganization(db, sessionID, orgName)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if orgID == 0 {
		t.Errorf("expected valid organization ID, got %d", orgID)
	}

	// Test duplicate organization name
	_, err = CreateOrganization(db, sessionID, orgName)
	if err == nil {
		t.Errorf("expected error for duplicate organization name, got nil")
	}
}
