package backend

import (
    "database/sql"
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

    // Create the necessary tables (USER, SESSION) for the test.
    _, err = db.Exec(`
        CREATE TABLE USER (
            id INTEGER PRIMARY KEY,
            name TEXT NOT NULL
        );
        CREATE TABLE SESSION (
            id INTEGER PRIMARY KEY,
            userid INTEGER NOT NULL
        );

        -- Insert sample data:
        INSERT INTO USER (id, name) VALUES (1, 'Alice');
        INSERT INTO SESSION (id, userid) VALUES (1, 1);
    `)
    if err != nil {
        t.Fatalf("failed to create schema or insert test data: %v", err)
    }

    // Call the refactored function, passing the in-memory DB and sessionID=1.
    err = ChangeDisplayName(db, 1, "Ivan")
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