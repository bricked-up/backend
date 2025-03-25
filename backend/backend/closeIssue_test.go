package backend

import (
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

func TestCloseIssue(t *testing.T) {
	tests := []struct {
		name        string
		sessionID   string
		issueID     string
		wantErr     bool
		wantDeleted bool
	}{
		{
			name:        "User1 closes Issue 1 => success",
			sessionID:   "1",
			issueID:     "1",
			wantErr:     false,
			wantDeleted: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup in-memory database
			db, err := sql.Open("sqlite", ":memory:")
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()

			// Load schema
			initSQL, err := os.ReadFile("../sql/init.sql")
			if err != nil {
				t.Fatal(err)
			}
			if _, err := db.Exec(string(initSQL)); err != nil {
				t.Fatal(err)
			}

			// Load test data
			populateSQL, err := os.ReadFile("../sql/populate.sql")
			if err != nil {
				t.Fatal(err)
			}
			if _, err := db.Exec(string(populateSQL)); err != nil {
				t.Fatal(err)
			}

			// Call the function
			err = CloseIssue(db, tc.issueID, tc.sessionID)

			// Assert error
			if tc.wantErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check if issue was deleted
			var count int
			err = db.QueryRow("SELECT COUNT(*) FROM issue WHERE id = ?", tc.issueID).Scan(&count)
			if err != nil {
				t.Fatal(err)
			}
			if tc.wantDeleted && count != 0 {
				t.Errorf("Issue not deleted")
			}
			if !tc.wantDeleted && count == 0 {
				t.Errorf("Issue was deleted unexpectedly")
			}
		})
	}
}
