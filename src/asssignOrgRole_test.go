package backend

import (
	"database/sql"
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
		return nil, err
	}
	if _, err := db.Exec(string(initSQL)); err != nil {
		return nil, err
	}

	// Execute populate.sql to insert test data
	populateSQL, err := os.ReadFile("../sql/populate.sql")
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(string(populateSQL)); err != nil {
		return nil, err
	}

	// Return the initialized database
	return db, nil
}

// Helper function to check if the user data was correctly inserted
func checkTestData(db *sql.DB) error {
	var userCount int
	err := db.QueryRow("SELECT COUNT(*) FROM USER").Scan(&userCount)
	if err != nil {
		return err
	}

	var orgCount int
	err = db.QueryRow("SELECT COUNT(*) FROM ORGANIZATION").Scan(&orgCount)
	if err != nil {
		return err
	}

	var roleCount int
	err = db.QueryRow("SELECT COUNT(*) FROM ORG_ROLE").Scan(&roleCount)
	if err != nil {
		return err
	}

	return nil
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
	if err := checkTestData(db); err != nil {
		t.Fatalf("Error checking test data: %v", err)
	}

	// Fetch user ID
	var userID int
	err = db.QueryRow("SELECT id FROM USER WHERE email = ?", "user1@example.com").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to fetch user ID: %v", err)
	}

	// Call the function under test with updated arguments
	// Assuming "Admin" role has an ID of 1
	err = assignOrgRole(db, userID, 1, 1, 1)

	// Assert no error occurred
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}
