package backend

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

// Helper function to initialize the in-memory database schema
func initializeTestDB(t *testing.T) (*sql.DB, error) {
	// Create an in-memory SQLite database using modernc driver
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}

	// Execute init.sql to create schema
	initSQL, err := os.ReadFile("../sql/init.sql")
	if err != nil {
		t.Fatalf("Failed to open init.sql: %v", err)
	}
	if _, err := db.Exec(string(initSQL)); err != nil {
		t.Fatalf("Failed to execute init.sql: %v", err)
	}

	// Execute populate.sql to insert test data
	populateSQL, err := os.ReadFile("../sql/populate.sql")
	if err != nil {
		t.Fatalf("Failed to open populate.sql: %v", err)
	}
	if _, err := db.Exec(string(populateSQL)); err != nil {
		t.Fatalf("Failed to execute populate.sql: %v", err)
	}

	// Return the initialized database
	return db, nil
}

// Helper function to check if the user data was correctly inserted
func checkTestData(db *sql.DB) {
	var userCount int
	err := db.QueryRow("SELECT COUNT(*) FROM USER").Scan(&userCount)
	if err != nil {
		fmt.Println("Error querying USER table:", err)
	} else {
		fmt.Printf("User count: %d\n", userCount)
	}

	var orgCount int
	err = db.QueryRow("SELECT COUNT(*) FROM ORGANIZATION").Scan(&orgCount)
	if err != nil {
		fmt.Println("Error querying ORGANIZATION table:", err)
	} else {
		fmt.Printf("Organization count: %d\n", orgCount)
	}

	var roleCount int
	err = db.QueryRow("SELECT COUNT(*) FROM ORG_ROLE").Scan(&roleCount)
	if err != nil {
		fmt.Println("Error querying ORG_ROLE table:", err)
	} else {
		fmt.Printf("Role count: %d\n", roleCount)
	}
}

// Test for the case where everything works fine
func TestAssignOrgRole_NoError(t *testing.T) {
	// Initialize in-memory SQLite database using init.sql and populate.sql
	db, err := initializeTestDB(t)
	if err != nil {
		t.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	// Check if data is correctly inserted
	checkTestData(db)

	// Call the function under test with updated arguments
	err = assignOrgRole(db, "user1@example.com", "Admin", 1, 1)

	// Assert no error occurred
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}
