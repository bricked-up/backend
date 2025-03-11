package backend

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

// TestDeleteUser tests the deleteUser function
func TestDeleteUser(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:") // Use in-memory database for testing
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			email TEXT,
			username TEXT,
			password TEXT
		);
		CREATE TABLE session (
			id TEXT PRIMARY KEY,
			user_id INTEGER
		);
	`)
	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	// Insert test user and session
	_, err = db.Exec("INSERT INTO users (id, email, username, password) VALUES (1, 'test@example.com', 'testuser', 'password123')")
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	_, err = db.Exec("INSERT INTO session (id, user_id) VALUES ('session123', 1)")
	if err != nil {
		t.Fatalf("Failed to insert test session: %v", err)
	}

	// Call deleteUser function
	err = deleteUser(db, "session123")
	if err != nil {
		t.Errorf("deleteUser returned an error: %v", err)
	}

	// Verify user deletion
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE id = 1").Scan(&count)
	if err != nil {
		t.Errorf("Failed to check user existence: %v", err)
	}

	if count != 0 {
		t.Errorf("User was not deleted, count: %d", count)
	}
}
