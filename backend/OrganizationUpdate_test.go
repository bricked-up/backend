package backend

import (
	"database/sql"
	"os"
	"testing"
)

func TestUpdateOrganizationName(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
	}

	// Load schema from init.sql
	initSQL, err := os.ReadFile("../sql/init.sql")
	if err != nil {
		t.Fatalf("failed to read init.sql: %v", err)
	}
	if _, err := db.Exec(string(initSQL)); err != nil {
		t.Fatalf("failed to execute init.sql: %v", err)
	}

	// Insert minimal test data
	_, err = db.Exec(`
		INSERT INTO ORGANIZATION (id, name) VALUES (1, 'TechCorp');
		INSERT INTO ORGANIZATION (id, name) VALUES (2, 'ExistingOrg');
		INSERT INTO USER (id, email, password, name, verified) VALUES (1, 'user1@example.com', 'pass', 'User One', 1);
		INSERT INTO SESSION (id, userid, expires) VALUES (1, 1, '2030-01-01 00:00:00');
	`)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Run a simple test
	err = UpdateOrganizationName(db, 1, 1, "New TechCorp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
