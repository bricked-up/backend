package projects

import (
	"brickedup/backend/utils"
	"encoding/json"
	"testing"

	_ "modernc.org/sqlite"
)

func TestGetTagDetails(t *testing.T) {
	db := utils.SetupTest(t)
	defer db.Close()

	// Test: Valid tag
	jsonData, err := getTagDetails(db, "1")
	if err != nil {
		t.Errorf("Expected valid tag, got error: %v", err)
	}

	var tag Tag
	err = json.Unmarshal(jsonData, &tag)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	if tag.ID != 1 {
		t.Errorf("Expected tag ID 1, got %d", tag.ID)
	}

	// Test: Invalid tag ID format
	_, err = getTagDetails(db, "abc")
	if err == nil || err.Error() != "invalid tag ID" {
		t.Errorf("Expected 'invalid tag ID' error, got: %v", err)
	}

	// Test: Non-existent tag
	_, err = getTagDetails(db, "999")
	if err == nil || err.Error() != "tag not found" {
		t.Errorf("Expected 'tag not found' error, got: %v", err)
	}
}
