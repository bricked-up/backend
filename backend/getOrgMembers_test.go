package backend

import (
	"database/sql"
	"encoding/json"
	"os"
	"reflect"
	"sort"
	"testing"

	_ "modernc.org/sqlite"
)

// setupTestDB creates an in-memory database and initializes it with the schema and test data
func setupTes(t *testing.T) *sql.DB {
	// Open an in-memory database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	// Read schema initialization SQL
	initSQL, err := os.ReadFile("../sql/init.sql")
	if err != nil {
		t.Fatalf("Failed to read init.sql: %v", err)
	}

	// Execute schema initialization
	_, err = db.Exec(string(initSQL))
	if err != nil {
		t.Fatalf("Failed to initialize database schema: %v", err)
	}

	// Read data population SQL
	populateSQL, err := os.ReadFile("../sql/populate.sql")
	if err != nil {
		t.Fatalf("Failed to read populate.sql: %v", err)
	}

	// Execute data population
	_, err = db.Exec(string(populateSQL))
	if err != nil {
		t.Fatalf("Failed to populate database: %v", err)
	}

	return db
}

func TestGetOrgMembers(t *testing.T) {
	// Setup test database
	db := setupTes(t)
	defer db.Close() // Ensure db is closed after all subtests complete

	// Test organization 1
	t.Run("Organization 1", func(t *testing.T) {
		orgID := 1

		// Get expected members directly from the database
		query := "SELECT userid FROM ORG_MEMBER WHERE orgid = ?"
		rows, err := db.Query(query, orgID)
		if err != nil {
			t.Fatalf("Failed to query expected members: %v", err)
		}

		var expectedMembers []int
		for rows.Next() {
			var userID int
			if err := rows.Scan(&userID); err != nil {
				t.Fatalf("Failed to scan row: %v", err)
			}
			expectedMembers = append(expectedMembers, userID)
		}
		rows.Close()

		// Call the function being tested
		result, err := GetOrgMembers(db, orgID)
		if err != nil {
			t.Fatalf("GetOrgMembers failed: %v", err)
		}

		// Parse the JSON result
		var actualMembers []int
		err = json.Unmarshal([]byte(result), &actualMembers)
		if err != nil {
			t.Fatalf("Failed to unmarshal result: %v", err)
		}

		// Sort both arrays for comparison
		sort.Ints(expectedMembers)
		sort.Ints(actualMembers)

		// Compare results
		if !reflect.DeepEqual(actualMembers, expectedMembers) {
			t.Errorf("Expected members %v, got %v", expectedMembers, actualMembers)
		}
	})

	// Test organization 2
	t.Run("Organization 2", func(t *testing.T) {
		orgID := 2

		// Get expected members directly from the database
		query := "SELECT userid FROM ORG_MEMBER WHERE orgid = ?"
		rows, err := db.Query(query, orgID)
		if err != nil {
			t.Fatalf("Failed to query expected members: %v", err)
		}

		var expectedMembers []int
		for rows.Next() {
			var userID int
			if err := rows.Scan(&userID); err != nil {
				t.Fatalf("Failed to scan row: %v", err)
			}
			expectedMembers = append(expectedMembers, userID)
		}
		rows.Close()

		// Call the function being tested
		result, err := GetOrgMembers(db, orgID)
		if err != nil {
			t.Fatalf("GetOrgMembers failed: %v", err)
		}

		// Parse the JSON result
		var actualMembers []int
		err = json.Unmarshal([]byte(result), &actualMembers)
		if err != nil {
			t.Fatalf("Failed to unmarshal result: %v", err)
		}

		// Sort both arrays for comparison
		sort.Ints(expectedMembers)
		sort.Ints(actualMembers)

		// Compare results
		if !reflect.DeepEqual(actualMembers, expectedMembers) {
			t.Errorf("Expected members %v, got %v", expectedMembers, actualMembers)
		}
	})

	// Test organization 3
	t.Run("Organization 3", func(t *testing.T) {
		orgID := 3

		// Get expected members directly from the database
		query := "SELECT userid FROM ORG_MEMBER WHERE orgid = ?"
		rows, err := db.Query(query, orgID)
		if err != nil {
			t.Fatalf("Failed to query expected members: %v", err)
		}

		var expectedMembers []int
		for rows.Next() {
			var userID int
			if err := rows.Scan(&userID); err != nil {
				t.Fatalf("Failed to scan row: %v", err)
			}
			expectedMembers = append(expectedMembers, userID)
		}
		rows.Close()

		// Call the function being tested
		result, err := GetOrgMembers(db, orgID)
		if err != nil {
			t.Fatalf("GetOrgMembers failed: %v", err)
		}

		// Parse the JSON result
		var actualMembers []int
		err = json.Unmarshal([]byte(result), &actualMembers)
		if err != nil {
			t.Fatalf("Failed to unmarshal result: %v", err)
		}

		// Sort both arrays for comparison
		sort.Ints(expectedMembers)
		sort.Ints(actualMembers)

		// Compare results
		if !reflect.DeepEqual(actualMembers, expectedMembers) {
			t.Errorf("Expected members %v, got %v", expectedMembers, actualMembers)
		}
	})

	// Test a non-existent organization
	t.Run("Non-existent Organization", func(t *testing.T) {
		nonExistentOrgID := 999

		result, err := GetOrgMembers(db, nonExistentOrgID)
		if err != nil {
			t.Fatalf("GetOrgMembers failed: %v", err)
		}

		var actualMembers []int
		err = json.Unmarshal([]byte(result), &actualMembers)
		if err != nil {
			t.Fatalf("Failed to unmarshal result: %v", err)
		}

		if len(actualMembers) != 0 {
			t.Errorf("Expected empty result for non-existent organization, got %v", actualMembers)
		}
	})
}

// Test that verifies the function handles database errors
func TestGetOrgMembersDBError(t *testing.T) {
	// Setup test database
	db := setupTes(t)

	// Create a separate test with its own db connection that we'll close
	closedDB := setupTes(t)
	closedDB.Close()

	// Attempt to get members from the closed DB, which should fail
	_, err := GetOrgMembers(closedDB, 1)
	if err == nil {
		t.Errorf("Expected a database error but got none")
	}

	// Clean up
	db.Close()
}
