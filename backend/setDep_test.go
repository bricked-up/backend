package backend

import (
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

// TestDeleteTag demonstrates using an in-memory DB to test the DeleteTag function.
// In this test the DB is opened in the testing for loop so it can properly generate
// a fresh db that corresponds to the expected effect of the function.
func TestSetDep(t *testing.T) {
	tests := []struct {
		name         string
		issueAid     int
		issueBid     int
		userid       int
		wantErr      bool
		wantInserted bool
	}{
		{
			name:         "Already Inserted dependency relation",
			issueBid:     3,
			issueAid:     1,
			userid:       1,
			wantErr:      true,
			wantInserted: false,
		},
		{
			name:         "User lacks permission",
			issueBid:     4,
			issueAid:     2,
			userid:       4, // assuming user 4 lacks write permission
			wantErr:      true,
			wantInserted: false,
		},
		{
			name:         "Invalid issue ID",
			issueBid:     1,
			issueAid:     999, // invalid issue ID
			userid:       1,
			wantErr:      true,
			wantInserted: false,
		},
		{
			name:         "Valid dependecy relation",
			issueBid:     5,
			issueAid:     3,
			userid:       3,
			wantErr:      false,
			wantInserted: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			//Open DB inside
			db, err := sql.Open("sqlite", ":memory:")
			if err != nil {
				t.Fatalf("failed to open in-memory db: %v", err)
			}
			defer db.Close()
			//Exec the init.sql
			initSQL, err := os.ReadFile("../sql/init.sql")
			if err != nil {
				t.Fatalf("failed to open init.sql: %v", err)
			}
			if _, err := db.Exec(string(initSQL)); err != nil {
				t.Fatalf("failed to exec init.sql: %v", err)
			}
			//Exec populate.sql
			populateSQL, err := os.ReadFile("../sql/populate.sql")
			if err != nil {
				t.Fatalf("failed to open populate.sql: %v", err)
			}
			if _, err := db.Exec(string(populateSQL)); err != nil {
				t.Fatalf("failed to exec populate.sql: %v", err)
			}
			// Call SetDep
			err = SetDep(db, tc.issueAid, tc.issueBid, tc.userid)
			if tc.wantErr && err == nil {
				t.Errorf("expected an error but got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("did NOT expect an error but got: %v", err)
			}
			// if we wnated some value inserted check if it was indeed inserted
			if tc.wantInserted {
				var count int
				err = db.QueryRow("SELECT COUNT(*) FROM DEPENDENCY WHERE issueid = ? AND dependency = ?", tc.issueBid, tc.issueAid).Scan(&count)
				if err != nil {
					t.Fatalf("failed to query DEPENDENCY table: %v", err)
				}
				if count != 1 {
					t.Errorf("expected dependency to be inserted, but it was not")
				}
			}
		})
	}
}
