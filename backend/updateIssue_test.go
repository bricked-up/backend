package backend

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func TestUpdateIssueDetails(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Load init.sql
	initSQL, err := os.ReadFile("../sql/init.sql")
	if err != nil {
		t.Fatalf("Failed to read init.sql: %v", err)
	}
	if _, err := db.Exec(string(initSQL)); err != nil {
		t.Fatalf("Failed to execute init.sql: %v", err)
	}

	// Load populate.sql
	populateSQL, err := os.ReadFile("../sql/populate.sql")
	if err != nil {
		t.Fatalf("Failed to read populate.sql: %v", err)
	}
	if _, err := db.Exec(string(populateSQL)); err != nil {
		t.Fatalf("Failed to execute populate.sql: %v", err)
	}

	// Insert session for user ID 1
	expiry := time.Now().Add(24 * time.Hour)
	res, err := db.Exec(`INSERT INTO SESSION (userid, expires) VALUES (?, ?)`, 1, expiry)
	if err != nil {
		t.Fatalf("Failed to insert session: %v", err)
	}
	sessionID, _ := res.LastInsertId()

	// Get issue ID from PROJECT_ISSUES table
	var issueID int
	err = db.QueryRow(`SELECT issueid FROM PROJECT_ISSUES WHERE id = 1`).Scan(&issueID)
	if err != nil {
		t.Fatalf("Failed to fetch issue ID: %v", err)
	}

	completedTime := time.Now()

	// Valid update test
	issue := Issue{
		Title:     "Updated Title",
		Desc:      "Updated description",
		Created:   time.Now(),
		Completed: &completedTime,
		Cost:      999,
	}
	err = UpdateIssueDetails(db, int(sessionID), issueID, issue)
	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}

	// Invalid issue ID test
	badIssue := Issue{
		Title:     "New Title",
		Desc:      "New Desc",
		Created:   time.Now(),
		Completed: nil,
		Cost:      3000,
	}
	err = UpdateIssueDetails(db, int(sessionID), 9999, badIssue)
	if err == nil {
		t.Error("Expected error for non-existent issue, got none")
	}
}
