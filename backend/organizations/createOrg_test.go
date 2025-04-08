package organizations

import (
	"brickedup/backend/utils"
	"testing"

	_ "modernc.org/sqlite"
)

func TestCreateOrganization(t *testing.T) {
	db := utils.SetupTest(t)
	defer db.Close()

	// Get an existing session ID from the database
	var sessionID int
	err := db.QueryRow("SELECT id FROM SESSION LIMIT 1").Scan(&sessionID)
	if err != nil {
		t.Fatalf("failed to get session ID: %v", err)
	}

	// Test valid organization creation
	orgName := "Test Organization Name"
	expectedSanitizedName := utils.SanitizeText(orgName, utils.TEXT)

	orgID, err := CreateOrganization(db, sessionID, orgName)
	if err != nil {
		t.Errorf("CreateOrganization returned error: %v", err)
	}

	if orgID == 0 {
		t.Errorf("expected valid organization ID, got %d", orgID)
	}

	// Verify organization was created with sanitized name
	var retrievedName string
	err = db.QueryRow("SELECT name FROM ORGANIZATION WHERE id = ?", orgID).Scan(&retrievedName)
	if err != nil {
		t.Errorf("failed to retrieve organization: %v", err)
	}
	if retrievedName != expectedSanitizedName {
		t.Errorf("expected organization name %s, got %s", expectedSanitizedName, retrievedName)
	}

	// Test duplicate organization name
	_, err = CreateOrganization(db, sessionID, orgName)
	if err == nil {
		t.Errorf("expected error for duplicate organization name, got nil")
	}

	// Test with potentially dangerous input
	maliciousName := "Dangerous'; DROP TABLE ORGANIZATION; --"
	sanitizedMalicious := utils.SanitizeText(maliciousName, utils.TEXT)

	if sanitizedMalicious == maliciousName {
		t.Errorf("sanitization failed to clean malicious input")
	}

	// Test with empty string after sanitization
	_, err = CreateOrganization(db, sessionID, "12345")
	if err == nil && utils.SanitizeText("12345", utils.TEXT) == "" {
		t.Errorf("should reject input that becomes empty after sanitization")
	}
}
