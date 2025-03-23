package backend

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func setupTeste(t *testing.T) *sql.DB {
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

func TestUpdateProject(t *testing.T) {
	db := setupTeste(t)
	defer db.Close()
	ctx := context.Background()

	// Get a valid project ID
	var projectID int
	err := db.QueryRow("SELECT id FROM PROJECT LIMIT 1").Scan(&projectID)
	if err != nil {
		t.Fatalf("Could not find a valid project: %v", err)
	}

	// Get original project data for comparison
	var originalProject Project
	err = db.QueryRow("SELECT id, orgid, name, budget, charter, archived FROM PROJECT WHERE id = ?", projectID).
		Scan(&originalProject.ID, &originalProject.OrgID, &originalProject.Name, &originalProject.Budget, &originalProject.Charter, &originalProject.Archived)
	if err != nil {
		t.Fatalf("Failed to get original project data: %v", err)
	}

	// Create a session with a valid user who has exec privileges
	// First find a user with exec privileges for this project
	var userID int
	err = db.QueryRow(`
		SELECT pm.userid FROM PROJECT_MEMBER pm
		JOIN PROJECT_MEMBER_ROLE pmr ON pm.id = pmr.memberid
		JOIN PROJECT_ROLE pr ON pmr.roleid = pr.id
		WHERE pm.projectid = ? AND pr.can_exec = 1
	`, projectID).Scan(&userID)
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

	// Create updated project data
	updatedProject := Project{
		ID:       originalProject.ID,
		OrgID:    originalProject.OrgID, // Keep the same org ID
		Name:     "Updated Project Name",
		Budget:   originalProject.Budget + 50000,
		Charter:  "Updated charter information for testing",
		Archived: !originalProject.Archived, // Toggle the archived status
	}

	// Test successful update
	err = updateProject(ctx, db, sessionID, projectID, updatedProject)
	if err != nil {
		t.Errorf("Expected project update to succeed, got error: %v", err)
	}

	// Verify the project was updated correctly
	var newProject Project
	err = db.QueryRow("SELECT id, orgid, name, budget, charter, archived FROM PROJECT WHERE id = ?", projectID).
		Scan(&newProject.ID, &newProject.OrgID, &newProject.Name, &newProject.Budget, &newProject.Charter, &newProject.Archived)
	if err != nil {
		t.Fatalf("Failed to get updated project data: %v", err)
	}

	// Compare updated fields
	if newProject.Name != updatedProject.Name {
		t.Errorf("Expected name to be '%s', got '%s'", updatedProject.Name, newProject.Name)
	}
	if newProject.Budget != updatedProject.Budget {
		t.Errorf("Expected budget to be %d, got %d", updatedProject.Budget, newProject.Budget)
	}
	if newProject.Charter != updatedProject.Charter {
		t.Errorf("Expected charter to be '%s', got '%s'", updatedProject.Charter, newProject.Charter)
	}
	if newProject.Archived != updatedProject.Archived {
		t.Errorf("Expected archived to be %v, got %v", updatedProject.Archived, newProject.Archived)
	}

	// Test with non-existent project ID
	nonExistentProjectID := 99999
	err = updateProject(ctx, db, sessionID, nonExistentProjectID, updatedProject)
	if err == nil {
		t.Errorf("Expected error for non-existent project ID, but got nil")
	}

	// Test with invalid session ID
	invalidSessionID := 99999
	err = updateProject(ctx, db, invalidSessionID, projectID, updatedProject)
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

	err = updateProject(ctx, db, expiredSessionID, projectID, updatedProject)
	if err == nil {
		t.Errorf("Expected error for expired session, but got nil")
	}

	// Test with user who doesn't have exec privileges
	// First we need to find or create a user without exec privileges
	var nonExecUserID int
	err = db.QueryRow(`
		SELECT u.id FROM USER u
		WHERE u.id NOT IN (
			SELECT pm.userid FROM PROJECT_MEMBER pm
			JOIN PROJECT_MEMBER_ROLE pmr ON pm.id = pmr.memberid
			JOIN PROJECT_ROLE pr ON pmr.roleid = pr.id
			WHERE pm.projectid = ? AND pr.can_exec = 1
		)
		LIMIT 1
	`, projectID).Scan(&nonExecUserID)

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

		// Add them to the project with a non-exec role
		memberResult, err := db.Exec("INSERT INTO PROJECT_MEMBER (userid, projectid) VALUES (?, ?)",
			nonExecUserID, projectID)
		if err != nil {
			t.Fatalf("Failed to create project member: %v", err)
		}
		memberID, err := memberResult.LastInsertId()
		if err != nil {
			t.Fatalf("Failed to get last insert ID: %v", err)
		}

		// Find a role without exec privileges
		var nonExecRoleID int
		err = db.QueryRow("SELECT id FROM PROJECT_ROLE WHERE can_exec = 0 LIMIT 1").Scan(&nonExecRoleID)
		if err != nil {
			// Create a non-exec role
			roleResult, err := db.Exec("INSERT INTO PROJECT_ROLE (name, can_view, can_edit, can_exec) VALUES (?, ?, ?, ?)",
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
		_, err = db.Exec("INSERT INTO PROJECT_MEMBER_ROLE (memberid, roleid) VALUES (?, ?)",
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
	err = updateProject(ctx, db, int(nonExecSessionID), projectID, updatedProject)
	if err == nil {
		t.Errorf("Expected error for user without exec privileges, but got nil")
	}
}
