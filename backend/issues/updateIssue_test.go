package issues

import (
	"brickedup/backend/utils"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func TestUpdateIssueDetails(t *testing.T) {
	db := utils.SetupTest(t)

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
