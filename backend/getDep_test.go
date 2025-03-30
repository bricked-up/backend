package backend

import (
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

func TestGetDep(t *testing.T) {
	// Open in-memory DB
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	defer db.Close()

	// Execute init.sql
	initSQL, err := os.ReadFile("../sql/init.sql")
	if err != nil {
		t.Fatalf("failed to open init.sql: %v", err)
	}
	if _, err := db.Exec(string(initSQL)); err != nil {
		t.Fatalf("failed to exec init.sql: %v", err)
	}

	// Execute populate.sql
	populateSQL, err := os.ReadFile("../sql/populate.sql")
	if err != nil {
		t.Fatalf("failed to open populate.sql: %v", err)
	}
	if _, err := db.Exec(string(populateSQL)); err != nil {
		t.Fatalf("failed to exec populate.sql: %v", err)
	}

	// Define subtests
	tests := []struct {
		name     string
		issueID  int
		wantErr  bool
		wantJSON string
	}{
		{
			name:     "Issue #3 has dependency 1",
			issueID:  3,
			wantErr:  false,
			wantJSON: `{"id":1}`,
		},
		{
			name:     "Issue #4 has dependency 2",
			issueID:  4,
			wantErr:  false,
			wantJSON: `{"id":2}`,
		},
		{
			name:     "Issue #999 - does not exist",
			issueID:  999,
			wantErr:  true,
			wantJSON: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := getDep(db, tc.issueID)

			if tc.wantErr && err == nil {
				t.Fatalf("expected an error but got none")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("did not expect an error but got: %v", err)
			}

			if err != nil {
				return
			}

			gotStr := string(got)
			if gotStr != tc.wantJSON {
				t.Errorf("JSON mismatch.\nGot:  %s\nWant: %s", gotStr, tc.wantJSON)
			}
		})
	}
}
