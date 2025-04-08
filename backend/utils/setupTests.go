package utils

import (
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

// SetupTest populates an in-memory Sqlite database for testing.
func SetupTest(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}
	initSQL, err := os.ReadFile("../../sql/init.sql")
	if err != nil {
		t.Fatalf("failed to read init.sql: %v", err)
	}
	_, err = db.Exec(string(initSQL))
	if err != nil {
		t.Fatalf("failed to execute init.sql: %v", err)
	}
	populateSQL, err := os.ReadFile("../../sql/populate.sql")
	if err != nil {
		t.Fatalf("failed to read populate.sql: %v", err)
	}
	_, err = db.Exec(string(populateSQL))
	if err != nil {
		t.Fatalf("failed to execute populate.sql: %v", err)
	}
	return db
}

