package projects

import (
	"brickedup/backend/utils"
	"reflect"
	"sort"
	"testing"

	_ "modernc.org/sqlite"
)

func TestGetProjMembers(t *testing.T) {
	// Setup test database
	db := utils.SetupTest(t)
	defer db.Close() // Ensure db is closed after all subtests complete

	// Test project 1
	t.Run("Organization 1", func(t *testing.T) {
		projectID := 1

		// Get expected members directly from the database
		query := "SELECT userid FROM PROJECT_MEMBER WHERE projectid = ?"
		rows, err := db.Query(query, projectID)
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
		result, err := getProjMembers(db, projectID)
		if err != nil {
			t.Fatalf("GetProjMembers failed: %v", err)
		}

		// Sort both arrays for comparison
		sort.Ints(expectedMembers)
		sort.Ints(result)

		// Compare results
		if !reflect.DeepEqual(result, expectedMembers) {
			t.Errorf("Expected members %v, got %v", expectedMembers, result)
		}
	})

	// Test a non-existent project
	t.Run("Non-existent Project", func(t *testing.T) {
		nonExistentProjID := 999

		result, err := getProjMembers(db, nonExistentProjID)
		if err != nil {
			t.Fatalf("GetProjMembers failed: %v", err)
		}

		if len(result) != 0 {
			t.Errorf("Expected empty result for non-existent organization, got %v", result)
		}
	})
}
