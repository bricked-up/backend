package projects

import (
	"brickedup/backend/utils"
	"testing"

	_ "modernc.org/sqlite"
)

func TestCreateTag(t *testing.T) {
	db := utils.SetupTest(t)
	defer db.Close()

	// Get an existing session ID and project ID from the database
	var sessionID, projectID int
	err := db.QueryRow("SELECT id FROM SESSION LIMIT 1").Scan(&sessionID)
	if err != nil {
		t.Fatalf("failed to get session ID: %v", err)
	}

	err = db.QueryRow("SELECT id FROM PROJECT LIMIT 1").Scan(&projectID)
	if err != nil {
		t.Fatalf("failed to get project ID: %v", err)
	}

	// Test valid tag creation
	tagName := "Test Tag Name"
	tagColor := "#FF5733" // A valid color code
	expectedSanitizedTagName := utils.SanitizeText(tagName, utils.TEXT)

	tagID, err := CreateTag(db, sessionID, projectID, tagName, tagColor)
	if err != nil {
		t.Errorf("CreateTag returned error: %v", err)
	}

	if tagID == 0 {
		t.Errorf("expected valid tag ID, got %d", tagID)
	}

	// Verify tag was created with sanitized name
	var retrievedTagName, retrievedTagColor string
	err = db.QueryRow("SELECT name, color FROM TAG WHERE id = ?", tagID).Scan(&retrievedTagName, &retrievedTagColor)
	if err != nil {
		t.Errorf("failed to retrieve tag: %v", err)
	}
	if retrievedTagName != expectedSanitizedTagName {
		t.Errorf("expected tag name %s, got %s", expectedSanitizedTagName, retrievedTagName)
	}
	if retrievedTagColor != tagColor {
		t.Errorf("expected tag color %s, got %s", tagColor, retrievedTagColor)
	}

	// Test duplicate tag name
	_, err = CreateTag(db, sessionID, projectID, tagName, tagColor)
	if err == nil {
		t.Errorf("expected error for duplicate tag name, got nil")
	}

	// Test with potentially dangerous input
	maliciousTagName := "Dangerous'; DROP TABLE TAG; --"
	sanitizedMalicious := utils.SanitizeText(maliciousTagName, utils.TEXT)

	if sanitizedMalicious == maliciousTagName {
		t.Errorf("sanitization failed to clean malicious input")
	}

	// Test with empty string after sanitization
	_, err = CreateTag(db, sessionID, projectID, "12345", "#000000")
	if err == nil && utils.SanitizeText("12345", utils.TEXT) == "" {
		t.Errorf("should reject input that becomes empty after sanitization")
	}

	// Test missing tagColor
	_, err = CreateTag(db, sessionID, projectID, tagName, "")
	if err == nil || err.Error() != "missing tagName or tagColor" {
		t.Errorf("expected error for missing tagColor, got: %v", err)
	}

	// Test missing tagName
	_, err = CreateTag(db, sessionID, projectID, "", tagColor)
	if err == nil || err.Error() != "missing tagName or tagColor" {
		t.Errorf("expected error for missing tagName, got: %v", err)
	}
}
