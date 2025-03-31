package backend

import (
	"database/sql"
	"encoding/json"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

func TestGetTagDetails(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Load schema
	initSQL, err := os.ReadFile("../sql/init.sql")
	if err != nil {
		t.Fatalf("Failed to read init.sql: %v", err)
	}
	if _, err := db.Exec(string(initSQL)); err != nil {
		t.Fatalf("Failed to execute init.sql: %v", err)
	}

	// Load data
	populateSQL, err := os.ReadFile("../sql/populate.sql")
	if err != nil {
		t.Fatalf("Failed to read populate.sql: %v", err)
	}
	if _, err := db.Exec(string(populateSQL)); err != nil {
		t.Fatalf("Failed to execute populate.sql: %v", err)
	}

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
