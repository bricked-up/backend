package backend

import (
	"database/sql"
	"os"
	"strconv"
	"testing"

	_ "modernc.org/sqlite"
)

// TestGetProjectDetails tests the GetProjectDetails function
func TestGetProjectDetails(t *testing.T) {
	// Initialize in-memory database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// Load and execute init.sql to set up the schema
	initSQL, err := os.ReadFile("../sql/init.sql")
	if err != nil {
		t.Fatalf("Failed to read init.sql: %v", err)
	}

	_, err = db.Exec(string(initSQL))
	if err != nil {
		t.Fatalf("Failed to execute init.sql: %v", err)
	}

	// Open and execute populate.sql to insert test data
	populateSQL, err := os.ReadFile("../sql/populate.sql")
	if err != nil {
		t.Fatalf("Failed to open populate.sql: %v", err)
	}

	_, err = db.Exec(string(populateSQL))
	if err != nil {
		t.Fatalf("Failed to execute populate.sql: %v", err)
	}

	// Define test cases
	tests := []struct {
		projectID int
		wantErr   bool
	}{
		{projectID: 1, wantErr: false},  // Assuming project ID 1 exists
		{projectID: 999, wantErr: true}, // Assuming project ID 999 does not exist
	}

	// Run tests
	for _, tt := range tests {
		t.Run("Testing project ID "+strconv.Itoa(tt.projectID), func(t *testing.T) {
			_, err := GetProjectDetails(db, tt.projectID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProjectDetails() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	// Close the database after tests are done
	defer db.Close()
}
