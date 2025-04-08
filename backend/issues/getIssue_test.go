package issues

import (
	"brickedup/backend/utils"
	"strconv"
	"testing"

	_ "modernc.org/sqlite"
)

// TestGetIssue tests the getIssue function
func TestGetIssue(t *testing.T) {
	db := utils.SetupTest(t)
	defer db.Close()

	tests := []struct {
		issueID int
		wantErr bool
	}{
		{issueID: 1, wantErr: false},  // Assuming issue ID 1 exists
		{issueID: 999, wantErr: true}, // Assuming issue ID 999 does not exist
	}

	for _, tt := range tests {
		t.Run("Testing issue ID "+strconv.Itoa(tt.issueID), func(t *testing.T) {
			_, err := GetIssue(db, tt.issueID)
			if (err != nil) != tt.wantErr {
				t.Errorf("getIssueDetails() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
