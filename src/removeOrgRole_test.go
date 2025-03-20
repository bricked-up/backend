package backend

import (
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTester(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}
	initSQL, err := os.ReadFile("../sql/init.sql")
	if err != nil {
		t.Fatalf("failed to read init.sql: %v", err)
	}
	_, err = db.Exec(string(initSQL))
	if err != nil {
		t.Fatalf("failed to execute init.sql: %v", err)
	}
	populateSQL, err := os.ReadFile("../sql/populate.sql")
	if err != nil {
		t.Fatalf("failed to read populate.sql: %v", err)
	}
	_, err = db.Exec(string(populateSQL))
	if err != nil {
		t.Fatalf("failed to execute populate.sql: %v", err)
	}
	return db
}

func TestRemoveOrgMemberRole(t *testing.T) {
	db := setupTester(t)
	defer db.Close()

	// Remove the role of user B in org (Developer role)
	var roleIDToRemove int
	err := db.QueryRow(`
        SELECT omr.id
        FROM ORG_MEMBER_ROLE omr
        JOIN ORG_MEMBER om ON omr.memberid = om.id
        JOIN ORG_ROLE r ON omr.roleid = r.id
        WHERE om.userid = 2 AND r.name = 'Developer'
    `).Scan(&roleIDToRemove)
	if err != nil {
		t.Fatalf("Could not find role to remove: %v", err)
	}

	// Valid session ID for user A
	validSessionID := 1

	// Verify the role exists before removal
	var countBefore int
	err = db.QueryRow("SELECT COUNT(*) FROM ORG_MEMBER_ROLE WHERE id = ?", roleIDToRemove).Scan(&countBefore)
	if err != nil {
		t.Fatalf("Failed to get initial count: %v", err)
	}
	if countBefore != 1 {
		t.Fatalf("Expected role to exist before removal, got count: %d", countBefore)
	}

	// Try to remove the role
	err = RemoveOrgMemberRole(db, validSessionID, roleIDToRemove)
	if err != nil {
		t.Errorf("Expected role removal to succeed, got error: %v", err)
	}

	// Verify the role is actually removed
	var countAfter int
	err = db.QueryRow("SELECT COUNT(*) FROM ORG_MEMBER_ROLE WHERE id = ?", roleIDToRemove).Scan(&countAfter)
	if err != nil {
		t.Errorf("Failed to check role existence: %v", err)
	}
	if countAfter != 0 {
		t.Errorf("Role should be removed, but count is %d", countAfter)
	}

	// Test with an invalid session ID
	invalidSessionID := 9999
	err = RemoveOrgMemberRole(db, invalidSessionID, roleIDToRemove)
	if err == nil {
		t.Errorf("Expected error for invalid session ID, but got nil")
	}
}
