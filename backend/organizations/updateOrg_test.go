package organizations

import (
	"brickedup/backend/utils"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func TestUpdateOrg(t *testing.T) {
	db := utils.SetupTest(t)
	defer db.Close()

	// Get a valid org ID
	var orgId int
	err := db.QueryRow("SELECT id FROM ORGANIZATION LIMIT 1").Scan(&orgId)
	if err != nil {
		t.Fatalf("Could not find a valid organization: %v", err)
	}

	// Get original organization data for comparison
	var originalOrg utils.Organization
	err = db.QueryRow("SELECT name FROM ORGANIZATION WHERE id = ?", orgId).
		Scan(&originalOrg.Name)
	if err != nil {
		t.Fatalf("Failed to get original organization data: %v", err)
	}

	// Create a session with a valid user who has exec privileges
	// First find a user with exec privileges for this project
	var userID int
	err = db.QueryRow(`
		SELECT om.userid FROM ORG_MEMBER om
		JOIN ORG_MEMBER_ROLE omr ON om.id = omr.memberid
		JOIN ORG_ROLE orgr ON omr.roleid = orgr.id
		WHERE om.orgid = ? AND orgr.can_exec = 1
	`, orgId).Scan(&userID)
	if err != nil {
		t.Fatalf("Could not find a user with exec privileges: %v", err)
	}

	// Create a session for this user
	result, err := db.Exec("INSERT INTO SESSION (userid, expires) VALUES (?, ?)",
		userID, time.Now().Add(1*time.Hour))
	if err != nil {
		t.Fatalf("Failed to create test session: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get last insert ID: %v", err)
	}
	sessionID := int(id)

	// Create updated org data
	updatedOrg := utils.Organization{
		ID:       originalOrg.ID,
		Name:     "Updated Org Name",
	}

	// Test successful update
	err = updateOrg(db, sessionID, orgId, updatedOrg)
	if err != nil {
		t.Errorf("Expected organization update to succeed, got error: %v", err)
	}

	// Verify the project was updated correctly
	var newOrg utils.Organization
	err = db.QueryRow("SELECT name FROM ORGANIZATION WHERE id = ?", orgId).
		Scan(&newOrg.Name)
	if err != nil {
		t.Fatalf("Failed to get updated project data: %v", err)
	}

	// Compare updated fields
	if newOrg.Name != updatedOrg.Name {
		t.Errorf("Expected name to be '%s', got '%s'", updatedOrg.Name, newOrg.Name)
	}

	// Test with non-existent org ID
	nonExistentOrgID := 99999
	err = updateOrg(db, sessionID, nonExistentOrgID, updatedOrg)
	if err == nil {
		t.Errorf("Expected error for non-existent org ID, but got nil")
	}

	// Test with invalid session ID
	invalidSessionID := 99999
	err = updateOrg(db, invalidSessionID, orgId, updatedOrg)
	if err == nil {
		t.Errorf("Expected error for invalid session ID, but got nil")
	}

	// Test with expired session
	// Create an expired session for testing
	expiredResult, err := db.Exec("INSERT INTO SESSION (userid, expires) VALUES (?, ?)",
		userID, time.Now().Add(-1*time.Hour))
	if err != nil {
		t.Fatalf("Failed to create expired test session: %v", err)
	}
	expiredID, err := expiredResult.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get last insert ID: %v", err)
	}
	expiredSessionID := int(expiredID)

	err = updateOrg(db, expiredSessionID, orgId, updatedOrg)
	if err == nil {
		t.Errorf("Expected error for expired session, but got nil")
	}

	// Test with user who doesn't have exec privileges
	// First we need to find or create a user without exec privileges
	var nonExecUserID int
	err = db.QueryRow(`
		SELECT u.id FROM USER u
		WHERE u.id NOT IN (
			SELECT om.userid FROM ORG_MEMBER om
			JOIN ORG_MEMBER_ROLE omr ON om.id = omr.memberid
			JOIN ORG_ROLE orgr ON omr.roleid = orgr.id
			WHERE om.orgid = ? AND orgr.can_exec = 1
		)
		LIMIT 1
	`, orgId).Scan(&nonExecUserID)

	if err != nil {
		// Create a new user
		userResult, err := db.Exec("INSERT INTO USER (name, email) VALUES (?, ?)",
			"Test User Without Privileges", "testuser@example.com")
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
		newUserID, err := userResult.LastInsertId()
		if err != nil {
			t.Fatalf("Failed to get last insert ID: %v", err)
		}
		nonExecUserID = int(newUserID)

		// Add them to the org with a non-exec role
		memberResult, err := db.Exec("INSERT INTO ORG_MEMBER (userid, orgid) VALUES (?, ?)",
			nonExecUserID, orgId)
		if err != nil {
			t.Fatalf("Failed to create org member: %v", err)
		}
		memberID, err := memberResult.LastInsertId()
		if err != nil {
			t.Fatalf("Failed to get last insert ID: %v", err)
		}

		// Find a role without exec privileges
		var nonExecRoleID int
		err = db.QueryRow("SELECT id FROM ORG_ROLE WHERE can_exec = 0 LIMIT 1").Scan(&nonExecRoleID)
		if err != nil {
			// Create a non-exec role
			roleResult, err := db.Exec("INSERT INTO ORG_ROLE (name, can_view, can_edit, can_exec) VALUES (?, ?, ?, ?)",
				"Viewer", 1, 0, 0)
			if err != nil {
				t.Fatalf("Failed to create viewer role: %v", err)
			}
			roleID, err := roleResult.LastInsertId()
			if err != nil {
				t.Fatalf("Failed to get last insert ID: %v", err)
			}
			nonExecRoleID = int(roleID)
		}

		// Assign the non-exec role
		_, err = db.Exec("INSERT INTO ORG_MEMBER_ROLE (memberid, roleid) VALUES (?, ?)",
			memberID, nonExecRoleID)
		if err != nil {
			t.Fatalf("Failed to assign role: %v", err)
		}
	}

	// Create a session for this non-exec user
	nonExecResult, err := db.Exec("INSERT INTO SESSION (userid, expires) VALUES (?, ?)",
		nonExecUserID, time.Now().Add(1*time.Hour))
	if err != nil {
		t.Fatalf("Failed to create test session: %v", err)
	}
	nonExecSessionID, err := nonExecResult.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get last insert ID: %v", err)
	}

	// Try to update with a user who doesn't have exec privileges
	err = updateOrg(db, int(nonExecSessionID), orgId, updatedOrg)
	if err == nil {
		t.Errorf("Expected error for user without exec privileges, but got nil")
	}
}
