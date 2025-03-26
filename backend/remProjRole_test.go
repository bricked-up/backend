package backend

import (
	"database/sql"
	"os"
	"testing"

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

	// Insert test users and setup
	_, err = db.Exec(`
		INSERT INTO USER (id, verifyid, email, password, name, verified) VALUES
		(10, 1, 'a@example.com', 'pw', 'User A', 1),
		(20, 1, 'b@example.com', 'pw', 'User B', 1),
		(30, 1, 'c@example.com', 'pw', 'User C', 1);

		INSERT INTO ORGANIZATION (id, name) VALUES (99, 'TestOrg');

		INSERT INTO PROJECT (id, orgid, name, budget, charter, archived) VALUES
		(99, 99, 'Test Project', 1000, 'Test Charter', 0);

		INSERT INTO PROJECT_ROLE (id, projectid, name, can_read, can_write, can_exec) VALUES
		(9901, 99, 'Manager', 1, 1, 1),
		(9902, 99, 'Developer', 1, 1, 0);

		INSERT INTO PROJECT_MEMBER (id, userid, projectid) VALUES
		(100, 10, 99),
		(200, 20, 99),
		(300, 30, 99);

		INSERT INTO PROJECT_MEMBER_ROLE (memberid, roleid) VALUES
		(100, 9901),
		(200, 9902),
		(300, 9902);
	`)
	if err != nil {
		t.Fatalf("Failed to insert test setup: %v", err)
	}

	// Insert session for user A (ID 10)
	res, err := db.Exec(`INSERT INTO SESSION (userid, expires) VALUES (?, ?)`, 10, "2030-01-01 10:00:00")
	if err != nil {
		t.Fatalf("Failed to insert session for user 10: %v", err)
	}
	sessionID, _ := res.LastInsertId()

	// Valid removal
	err = removeUserRole(db, int(sessionID), "20", 9902, 99)
	if err != nil {
		t.Fatalf("Valid removal failed: %v", err)
	}

	// Confirm role was removed
	var exists bool
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM PROJECT_MEMBER_ROLE
			WHERE memberid = 200 AND roleid = 9902
		)
	`).Scan(&exists)
	if err != nil {
		t.Fatalf("Failed to query role: %v", err)
	}
	if exists {
		t.Fatal("Role was not removed as expected")
	}

	// User B does not exist
	err = removeUserRole(db, int(sessionID), "999", 9902, 99)
	if err == nil || err.Error() != "user B does not exist" {
		t.Fatalf("Expected 'user B does not exist', got: %v", err)
	}

	// Reassign role to user B before next test
	_, err = db.Exec(`INSERT INTO PROJECT_MEMBER_ROLE (memberid, roleid) VALUES (?, ?)`, 200, 9902)
	if err != nil {
		t.Fatalf("Failed to reassign role to user B: %v", err)
	}

	// User A without exec (user C)
	res, err = db.Exec(`INSERT INTO SESSION (userid, expires) VALUES (?, ?)`, 30, "2030-01-01 10:00:00")
	if err != nil {
		t.Fatalf("Failed to insert session for user 30: %v", err)
	}
	noExecID, _ := res.LastInsertId()

	err = removeUserRole(db, int(noExecID), "20", 9902, 99)
	if err == nil || err.Error() != "user A lacks exec permissions" {
		t.Fatalf("Expected 'user A lacks exec permissions', got: %v", err)
	}
}
