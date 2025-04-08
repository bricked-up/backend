package organizations

import (
	"brickedup/backend/utils"
	"encoding/json"
	"reflect"
	"sort"
	"testing"

	_ "modernc.org/sqlite"
)

func TestGetOrgMembers(t *testing.T) {
	// Setup test database
	db := utils.SetupTest(t)
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
