package backend

import (
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

// setupDatabase connects to the database and initializes it with init.sql and populate.sql data
func setupDatabase(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	initSQL, err := os.ReadFile("../sql/init.sql")
	if err != nil {
		t.Fatalf("Failed to read init.sql: %v", err)
	}

	_, err = db.Exec(string(initSQL))
	if err != nil {
		t.Fatalf("Failed to initialize database schema: %v", err)
	}

	populateSQL, err := os.ReadFile("../sql/populate.sql")
	if err != nil {
		t.Fatalf("Failed to read populate.sql: %v", err)
	}

	_, err = db.Exec(string(populateSQL))
	if err != nil {
		t.Fatalf("Failed to populate database: %v", err)
	}

	return db
}

func TestAssignIssuePriority(t *testing.T) {
	db := setupDatabase(t)
	defer db.Close()

	t.Run("Successful priority assignment", func(t *testing.T) {
		err := AssignIssuePriority(db, 1, 1, 2)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("Invalid session ID", func(t *testing.T) {
		err := AssignIssuePriority(db, 999, 1, 2)
		if err == nil || err.Error() != "Invalid session ID" {
			t.Errorf("Expected error 'Invalid session ID', got %v", err)
		}
	})

	t.Run("Invalid priority for project", func(t *testing.T) {
		err := AssignIssuePriority(db, 1, 1, 999)
		if err == nil || err.Error() != "Invalid priority for the project" {
			t.Errorf("Expected error 'Invalid priority for the project', got %v", err)
		}
	})

	t.Run("User without write permissions", func(t *testing.T) {
		_, _ = db.Exec("UPDATE PROJECT_ROLE SET can_write = 0 WHERE id = 1")
		err := AssignIssuePriority(db, 1, 1, 2)
		if err == nil || err.Error() != "User does not have write permissions" {
			t.Errorf("Expected error 'User does not have write permissions', got %v", err)
		}
	})
}
