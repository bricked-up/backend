package backend

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func TestRemoveUserRole(t *testing.T) {
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

	// Load base data
	populateSQL, err := os.ReadFile("../sql/populate.sql")
	if err != nil {
		t.Fatalf("Failed to read populate.sql: %v", err)
	}
	if _, err := db.Exec(string(populateSQL)); err != nil {
		t.Fatalf("Failed to execute populate.sql: %v", err)
	}

	// Insert session for user A (userID 1 - John, has exec permission in project 1)
	expiry := time.Now().Add(24 * time.Hour)
	res, err := db.Exec(`INSERT INTO SESSION (userid, expires) VALUES (?, ?)`, 1, expiry)
	if err != nil {
		t.Fatalf("Failed to insert session for user 1: %v", err)
	}
	sessionID, _ := res.LastInsertId()

	// SUCCESS: Remove Jane (userID 2)'s Developer role (roleid=2) in project 1
	err = removeUserRole(db, int(sessionID), 2, 2, 1)
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}

	// Confirm role was removed
	var exists bool
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM PROJECT_MEMBER_ROLE
			WHERE memberid = (SELECT id FROM PROJECT_MEMBER WHERE userid = 2 AND projectid = 1)
			AND roleid = 2
		)
	`).Scan(&exists)
	if err != nil {
		t.Fatalf("Failed to query role: %v", err)
	}
	if exists {
		t.Fatal("Role was not removed as expected")
	}

	// ERROR: User B does not exist
	err = removeUserRole(db, int(sessionID), 999, 2, 1)
	if err == nil || err.Error() != "user B does not exist" {
		t.Fatalf("Expected 'user B does not exist', got: %v", err)
	}

	// Reassign role to Jane before next test
	_, err = db.Exec(`
		INSERT INTO PROJECT_MEMBER_ROLE (memberid, roleid)
		VALUES (
			(SELECT id FROM PROJECT_MEMBER WHERE userid = 2 AND projectid = 1),
			2
		)
	`)
	if err != nil {
		t.Fatalf("Failed to reassign role: %v", err)
	}

	// ERROR: User A lacks exec permissions (userID=2, Jane)
	res, err = db.Exec(`INSERT INTO SESSION (userid, expires) VALUES (?, ?)`, 2, expiry)
	if err != nil {
		t.Fatalf("Failed to insert session for user 2: %v", err)
	}
	noExecSessionID, _ := res.LastInsertId()

	err = removeUserRole(db, int(noExecSessionID), 1, 1, 1)
	if err == nil || err.Error() != "user A lacks exec permissions" {
		t.Fatalf("Expected 'user A lacks exec permissions', got: %v", err)
	}
}
