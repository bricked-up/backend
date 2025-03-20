package backend

import (
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

func TestLogin(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:") // Use in-memory database
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Load database schema
	initSQL, err := os.ReadFile("../sql/init.sql")
	if err != nil {
		t.Fatalf("Failed to read init.sql: %v", err)
	}
	if _, err := db.Exec(string(initSQL)); err != nil {
		t.Fatalf("Failed to execute init.sql: %v", err)
	}

	// Load initial data
	populateSQL, err := os.ReadFile("../sql/populate.sql")
	if err != nil {
		t.Fatalf("Failed to read populate.sql: %v", err)
	}
	if _, err := db.Exec(string(populateSQL)); err != nil {
		t.Fatalf("Failed to execute populate.sql: %v", err)
	}


	// Test valid login
	sessionid, err := login(db, "john.doe@example.com", "hashed_password_1")
    if err != nil {
        t.Fatal(err)
    }

	if sessionid < 0 {
        t.Fatalf("Valid login failed: invalid sessionid")
    }

    // Test invalid password
    _, err = login(db, "user1@example.com", "wrongpassword")
    if err == nil {
        t.Fatal("Invalid password failed: should not be logged in")
    }

	// Test non-existent user
	_, err = login(db, "nouser@example.com", "testpassword")
    if err == nil {
        t.Fatal("Non-existent user failed: should not be logged in")
    }

	// Test unverified user
	_, err = login(db, "unverified@example.com", "password3")
    if err == nil {
        t.Fatal("Unverified user failed: should not be logged in")
    }
}
