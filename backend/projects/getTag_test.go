package projects

import (
	"brickedup/backend/utils"
	"testing"

	_ "modernc.org/sqlite"
)

func TestGetTagDetails(t *testing.T) {
	db := utils.SetupTest(t)
	defer db.Close()

	// Test: Valid tag
	_, err := GetTag(db, 1)
	if err != nil {
		t.Errorf("Expected valid tag, got error: %v", err)
	}

	// Test: Non-existent tag
	_, err = GetTag(db, 999)
	if err == nil || err.Error() != "tag not found" {
		t.Errorf("Expected 'tag not found' error, got: %v", err)
	}
}
