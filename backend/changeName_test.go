package backend

import (
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

// TestUpdateUserName demonstrates using an in-memory DB to test ChangeDisplayName
func TestUpdateUserName(t *testing.T) {
    // Open an in-memory database with the modernc driver name "sqlite".
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        t.Fatalf("failed to open in-memory db: %v", err)
    }
    defer db.Close()

	// Open the init.sql file and exec:
	initSQL, err := os.ReadFile("../sql/init.sql")
	if err != nil {
		t.Fatalf("failed to open init.sql %v", err)
	}
	defer db.Close()

	if _, err := db.Exec(string(initSQL)); err != nil {
        t.Fatalf("failed to exec init.sql: %v", err)
    }
	
	// Open the populate.sql file and exec:
	populateSQL, err := os.ReadFile("../sql/populate.sql")
	if err != nil {
		t.Fatalf("failed to open populate.sql %v", err)
	}
	defer db.Close()

	if _, err := db.Exec(string(populateSQL)); err != nil {
        t.Fatalf("failed to exec populate.sql: %v", err)
    }

    // Call the refactored function, passing the in-memory DB and sessionID=1.
    err = ChangeDisplayName(db, 1, "Ivan123!!")
    if err != nil {
        t.Errorf("ChangeDisplayName returned error: %v", err)
    }

    // Verify the name was updated to "Ivan".
    var updatedName string
    err = db.QueryRow("SELECT name FROM USER WHERE id = 1").Scan(&updatedName)
    if err != nil {
        t.Errorf("failed to query updated name: %v", err)
    }
    if updatedName != "Ivan" {
        t.Errorf("expected 'Ivan', got '%s'", updatedName)
    }
}