package backend

import (
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

func TestPromoteUserRole(t *testing.T) {
	// Setup database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Load schema
	initSQL, err := os.ReadFile("../sql/init.sql")
	if err != nil {
		t.Fatalf("Failed to read init.sql: %v", err)
	}
	_, err = db.Exec(string(initSQL))
	if err != nil {
		t.Fatalf("Failed to initialize schema: %v", err)
	}

	// Load data
	populateSQL, err := os.ReadFile("../sql/populate.sql")
	if err != nil {
		t.Fatalf("Failed to read populate.sql: %v", err)
	}
	_, err = db.Exec(string(populateSQL))
	if err != nil {
		t.Fatalf("Failed to populate database: %v", err)
	}

	// Test Case 1: Valid promotion
	err = promoteUserRole(db, 1, 2, 3, 1) // User 1 promotes User 2 to role 3 in project 1
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test Case 2: User B already has the specified role
	err = promoteUserRole(db, 1, 2, 3, 1)
	if err == nil || err.Error() != "user B already has the specified role" {
		t.Errorf("Expected 'user B already has the specified role', got %v", err)
	}

	// Test Case 3: User A lacks exec permissions (user 2 in project 1 has Developer role)
	err = promoteUserRole(db, 2, 1, 2, 1)
	if err == nil || err.Error() != "user A lacks exec permissions" {
		t.Errorf("Expected 'user A lacks exec permissions', got %v", err)
	}

	// Test Case 4: User B not part of the project (user 4 is not in project 1)
	err = promoteUserRole(db, 1, 4, 2, 1)
	if err == nil || err.Error() != "user B is not part of the project" {
		t.Errorf("Expected 'user B is not part of the project', got %v", err)
	}

	// Test Case 5: User B not validated (user 4 is unverified in populate.sql)
	_, err = db.Exec(`UPDATE USER SET verifyid = NULL WHERE id = 4`)
	if err != nil {
		t.Fatalf("Failed to modify verification status: %v", err)
	}
	err = promoteUserRole(db, 1, 4, 2, 4)
	if err == nil || err.Error() != "user B is not validated" {
		t.Errorf("Expected 'user B is not validated', got %v", err)
	}

	// Test Case 6: User B does not exist
	err = promoteUserRole(db, 1, 999, 2, 1)
	if err == nil || err.Error() != "user B does not exist" {
		t.Errorf("Expected 'user B does not exist', got %v", err)
	}
}
