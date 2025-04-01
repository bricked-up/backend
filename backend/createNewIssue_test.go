package backend

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func TestCreateNewIssue(t *testing.T) {
	// Open an in-memory SQLite database.
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	defer db.Close()

	// Load and execute the schema from init.sql
	initSQL, err := os.ReadFile("../sql/init.sql")
	if err != nil {
		t.Fatalf("failed to open init.sql: %v", err)
	}
	if _, err := db.Exec(string(initSQL)); err != nil {
		t.Fatalf("failed to execute init.sql: %v", err)
	}

	// Load and execute initial data from populate.sql
	populateSQL, err := os.ReadFile("../sql/populate.sql")
	if err != nil {
		t.Fatalf("failed to open populate.sql: %v", err)
	}
	if _, err := db.Exec(string(populateSQL)); err != nil {
		t.Fatalf("failed to execute populate.sql: %v", err)
	}

	// Test data

	projectid := 1;
	title := "Sample Issue"
	desc := "This is a sample issue"
	tagid := 1    // Ensure this matches an existing ID in TAG table
	priority := 1 // Ensure this matches an existing ID in PRIORITY table
	completed := time.Now().Add(24 * time.Hour)
	cost := 500
	createdDate := time.Now()

	// Create a session with a valid user who has exec privileges
	// First find a user with exec privileges for this project
	var userID int
	err = db.QueryRow(`
		SELECT pm.userid FROM PROJECT_MEMBER pm
		JOIN PROJECT_MEMBER_ROLE pmr ON pm.id = pmr.memberid
		JOIN PROJECT_ROLE pr ON pmr.roleid = pr.id
		WHERE pm.projectid = ? AND pr.can_write = 1
	`, projectid).Scan(&userID)
	if err != nil {
		t.Fatalf("Could not find a user with write privileges: %v", err)
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

	// Call the function to test
	_, err = CreateIssue(
		sessionID, 
		projectid, 
		title, 
		desc, 
		tagid, 
		priority, 
		completed, 
		cost, 
		createdDate, 
		db)

	if err != nil {
		t.Fatalf("Failed to create new issue: %v", err)
	}

	// Non-existant sessionid
	nonExistentSession := 999999
	_, err = CreateIssue(
		nonExistentSession, 
		projectid, 
		title, 
		desc, 
		tagid, 
		priority, 
		completed, 
		cost, 
		createdDate, 
		db)

	if err == nil {
		t.Fatalf("Invalid sessionid should fail.")
	}

	// Non-existant projct
	nonExistentProject := 999999
	_, err = CreateIssue(
		sessionID, 
		nonExistentProject, 
		title, 
		desc, 
		tagid, 
		priority, 
		completed, 
		cost, 
		createdDate, 
		db)

	if err == nil {
		t.Fatalf("Invalid projectid should fail.")
	}

	// Test with user who doesn't have exec privileges
	// First we need to find or create a user without exec privileges
	var nonWriteUserID int
	err = db.QueryRow(`
		SELECT u.id FROM USER u
		WHERE u.id NOT IN (
			SELECT pm.userid FROM PROJECT_MEMBER pm
			JOIN PROJECT_MEMBER_ROLE pmr ON pm.id = pmr.memberid
			JOIN PROJECT_ROLE pr ON pmr.roleid = pr.id
			WHERE pm.projectid = ? AND pr.can_write = 1
		)
		LIMIT 1
	`, projectid).Scan(&nonWriteUserID)

	if err != nil {
		t.Fatalf("User without write privileges does not exist in the project!")
	}

	_, err = CreateIssue(
		sessionID, 
		nonExistentProject, 
		title, 
		desc, 
		tagid, 
		priority, 
		completed, 
		cost, 
		createdDate, 
		db)

	if err == nil {
		t.Fatalf("Only write-allowed users should be able to create new issues!")
	}

}
