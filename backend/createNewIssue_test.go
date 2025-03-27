package backend

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func TestCreateNewIssue(t *testing.T) {
	// Open an in-memory SQLite database.
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	defer db.Close()

	// Load and execute the schema from init.sql
	initSQL, err := os.ReadFile("../sql/init.sql")
	if err != nil {
		t.Fatalf("failed to open init.sql: %v", err)
	}
	if _, err := db.Exec(string(initSQL)); err != nil {
		t.Fatalf("failed to execute init.sql: %v", err)
	}

	// Load and execute initial data from populate.sql
	populateSQL, err := os.ReadFile("../sql/populate.sql")
	if err != nil {
		t.Fatalf("failed to open populate.sql: %v", err)
	}
	if _, err := db.Exec(string(populateSQL)); err != nil {
		t.Fatalf("failed to execute populate.sql: %v", err)
	}

	// Test data

	title := "Sample Issue"
	desc := "This is a sample issue"
	tagid := 1      // Ensure this matches an existing ID in TAG table
	priority := 1 // Ensure this matches an existing ID in PRIORITY table
	completed := time.Now().Add(24 * time.Hour)
	cost := 500
	createdDate := time.Now()

	// Call the function to test
	insertedID, err := CreateIssue(title, desc, tagid, priority, completed, cost, createdDate, db)
	if err != nil {
		t.Fatalf("Failed to create new issue: %v", err)
	}

	// Validate the insertion by querying the database
	var retrievedTitle string
	err = db.QueryRow("SELECT title FROM issue WHERE id = ?", insertedID).Scan(&retrievedTitle)
	if err != nil {
		t.Fatalf("Failed to retrieve inserted issue: %v", err)
	}

	// Assert the retrieved title matches the input
	if retrievedTitle != title {
		t.Errorf("Expected title '%s', got '%s'", title, retrievedTitle)
	}
}
