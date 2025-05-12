package projects

import (
	"brickedup/backend/utils"
	"testing"

	_ "modernc.org/sqlite"
)

func TestCreateProj(t *testing.T) {
	db := utils.SetupTest(t)
	defer db.Close()

	// Get an existing session ID from the database
	var sessionID int
	err := db.QueryRow("SELECT id FROM SESSION LIMIT 1").Scan(&sessionID)
	if err != nil {
		t.Fatalf("failed to get session ID: %v", err)
	}

	projName := "Test Project Name"

	err = CreateProj(db, sessionID, 1, projName, 500000, "some charter")
	if err != nil {
		t.Errorf("CreateOrganization returned error: %v", err)
	}
}
