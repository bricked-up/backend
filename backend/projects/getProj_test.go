package projects

import (
	"brickedup/backend/utils"
	"strconv"
	"testing"

	_ "modernc.org/sqlite"
)

// TestGetProject tests the GetProjectDetails function
func TestGetProject(t *testing.T) {
	db := utils.SetupTest(t)
	defer db.Close()

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
			_, err := GetProject(db, tt.projectID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
