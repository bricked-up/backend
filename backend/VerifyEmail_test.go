package backend

import (
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTestDb(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	initSQL, err := os.ReadFile("../sql/init.sql")
	if err != nil {
		t.Fatalf("failed to read init.sql: %v", err)
	}
	if _, err := db.Exec(string(initSQL)); err != nil {
		t.Fatalf("failed to exec init.sql: %v", err)
	}

	populateSQL, err := os.ReadFile("../sql/populate.sql")
	if err != nil {
		t.Fatalf("failed to read populate.sql: %v", err)
	}
	if _, err := db.Exec(string(populateSQL)); err != nil {
		t.Fatalf("failed to exec populate.sql: %v", err)
	}

	return db
}

func TestVerifyUser_FromPopulateSQL(t *testing.T) {
	tests := []struct {
		name       string
		code       int
		expectErr  bool
		userID     int
		expectNull bool
		expectTrue bool
	}{
		{
			name:       "Valid code from populate.sql",
			code:       123456, // this exists and is not expired
			expectErr:  false,
			userID:     1,
			expectNull: true,
			expectTrue: true,
		},
		{
			name:       "Non-existent code",
			code:       999999,
			expectErr:  true,
			userID:     1,
			expectNull: false,
			expectTrue: true, // Already verified
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			defer db.Close()

			err := VerifyUser(tc.code, db)
			if tc.expectErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("did not expect error but got: %v", err)
			}

			var verifyID sql.NullInt64
			var verified bool
			err = db.QueryRow("SELECT verifyid, verified FROM USER WHERE id = ?", tc.userID).Scan(&verifyID, &verified)
			if err != nil {
				t.Fatalf("failed to query user: %v", err)
			}

			if tc.expectNull && verifyID.Valid {
				t.Errorf("expected verifyid to be NULL, got %v", verifyID.Int64)
			}
			if tc.expectTrue && !verified {
				t.Errorf("expected verified to be true, got false")
			}
		})
	}
}
