package backend

import (
	"database/sql"
	"testing"
	"os"

	_ "modernc.org/sqlite"
)


// TestDeleteTag demonstrates using an in-memory DB to test the DeleteTag function.
// In this test the DB is opened in the testing for loop so it can properly generate
// a fresh db that corresponds to the expected effect of the function.
func TestDeleteTag(t *testing.T) {
    tests := []struct {
        name        string
        sessionID   int
        tagID       int
        wantErr     bool
        wantDeleted bool
    }{
        {
            name:        "User1 can delete Tag #1 => success",
            sessionID:   1,
            tagID:       1,
            wantErr:     false,
            wantDeleted: true,
        },
        {
            name:        "User4 tries to delete Tag #1 => should fail",
            sessionID:   5,
            tagID:       1,
            wantErr:     true,
            wantDeleted: false,
        },
        {
            name:        "Invalid session => should fail",
            sessionID:   999,
            tagID:       1,
            wantErr:     true,
            wantDeleted: false,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            // Open a DB *inside* each subtest
            db, err := sql.Open("sqlite", ":memory:")
            if err != nil {
                t.Fatalf("failed to open in-memory db: %v", err)
            }
            defer db.Close()

            // Exec init.sql
            initSQL, err := os.ReadFile("../sql/init.sql")
            if err != nil {
                t.Fatalf("failed to open init.sql: %v", err)
            }
            if _, err := db.Exec(string(initSQL)); err != nil {
                t.Fatalf("failed to exec init.sql: %v", err)
            }

            // Exec populate.sql
            populateSQL, err := os.ReadFile("../sql/populate.sql")
            if err != nil {
                t.Fatalf("failed to open populate.sql: %v", err)
            }
            if _, err := db.Exec(string(populateSQL)); err != nil {
                t.Fatalf("failed to exec populate.sql: %v", err)
            }

            // Call DeleteTag
            err = DeleteTag(db, tc.sessionID, tc.tagID)
            if tc.wantErr && err == nil {
                t.Errorf("expected an error but got nil")
            }
            if !tc.wantErr && err != nil {
                t.Errorf("did NOT expect an error but got: %v", err)
            }

            // Check whether the tag was deleted or not
            var count int
            queryErr := db.QueryRow("SELECT COUNT(*) FROM TAG WHERE id = ?", tc.tagID).Scan(&count)
            if queryErr != nil {
                t.Fatalf("failed to query TAG table: %v", queryErr)
            }

            if tc.wantDeleted && count != 0 {
                t.Errorf("expected tag (id=%d) to be deleted, but it still exists", tc.tagID)
            } else if !tc.wantDeleted && count == 0 {
                t.Errorf("expected tag (id=%d) NOT to be deleted, but got count=0", tc.tagID)
            }
        })
    }
}
