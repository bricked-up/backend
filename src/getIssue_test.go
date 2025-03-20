package backend

import (
	"database/sql"
	"os"
	"strconv"
	"testing"

	_ "modernc.org/sqlite"
)

// setupTest initializes the in-memory database for testing
func setupTest(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	initSQL, err := os.ReadFile("../sql/init.sql")
	if err != nil {
		t.Fatalf("Failed to read init.sql: %v", err)
	}

	_, err = db.Exec(string(initSQL))
	if err != nil {
		t.Fatalf("Failed to execute init.sql: %v", err)
	}

	populateSQL, err := os.ReadFile("../sql/populate.sql")
	if err != nil {
		t.Fatalf("Failed to read populate.sql: %v", err)
	}

	_, err = db.Exec(string(populateSQL))
	if err != nil {
		t.Fatalf("Failed to execute populate.sql: %v", err)
	}

	return db
}

// TestGetIssueDetails tests the getIssueDetails function
func TestGetIssueDetails(t *testing.T) {
	db := setupTest(t)
	defer db.Close()

	tests := []struct {
		issueID int
		wantErr bool
	}{
		{issueID: 1, wantErr: false},  // Assuming issue ID 1 exists
		{issueID: 999, wantErr: true}, // Assuming issue ID 999 does not exist
	}

	for _, tt := range tests {
		t.Run("Testing issue ID "+strconv.Itoa(tt.issueID), func(t *testing.T) {
			_, err := getIssueDetails(db, tt.issueID)
			if (err != nil) != tt.wantErr {
				t.Errorf("getIssueDetails() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
