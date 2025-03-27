package backend

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func TestAssignProjectRole(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Load init.sql
	initSQL, err := os.ReadFile("../sql/init.sql")
	if err != nil {
		t.Fatalf("Failed to read init.sql: %v", err)
	}
	if _, err := db.Exec(string(initSQL)); err != nil {
		t.Fatalf("Failed to execute init.sql: %v", err)
	}

	// Load populate.sql
	populateSQL, err := os.ReadFile("../sql/populate.sql")
	if err != nil {
		t.Fatalf("Failed to read populate.sql: %v", err)
	}
	if _, err := db.Exec(string(populateSQL)); err != nil {
		t.Fatalf("Failed to execute populate.sql: %v", err)
	}

	// Insert fresh session for user ID 1 (John)
	expiry := time.Now().Add(24 * time.Hour)
	result, err := db.Exec(`INSERT INTO SESSION (userid, expires) VALUES (?, ?)`, 1, expiry)
	if err != nil {
		t.Fatalf("Failed to insert test session: %v", err)
	}
	sessionID, _ := result.LastInsertId()

	// SUCCESS: John (userID=1) promotes Jane (userID=2) to roleid=3 (QA Tester) in project 1
	err = assignProjectRole(db, int(sessionID), 2, 3, 1)
	if err != nil {
		t.Errorf("expected success but got error: %v", err)
	}

	// ERROR: Invalid session
	err = assignProjectRole(db, 9999, 2, 3, 1)
	if err == nil {
		t.Errorf("expected error for invalid session but got none")
	}

	// ERROR: Non-existent project
	err = assignProjectRole(db, int(sessionID), 2, 3, 999)
	if err == nil {
		t.Errorf("expected error for non-existent project but got none")
	}

	// ERROR: Unverified user (userID=4)
	err = assignProjectRole(db, int(sessionID), 4, 5, 4)
	if err == nil {
		t.Errorf("expected error for unverified user but got none")
	}

	// ERROR: User not in project (userID=4 not in project 1)
	err = assignProjectRole(db, int(sessionID), 4, 3, 1)
	if err == nil {
		t.Errorf("expected error for user not in project but got none")
	}

	// ERROR: User already has role (Jane already has roleid=2 in project 1)
	err = assignProjectRole(db, int(sessionID), 2, 2, 1)
	if err == nil {
		t.Errorf("expected error for duplicate role but got none")
	}

	// ERROR: User without exec permission (Jane (userID=2), sessionID needs fresh insert)
	result, err = db.Exec(`INSERT INTO SESSION (userid, expires) VALUES (?, ?)`, 2, expiry)
	if err != nil {
		t.Fatalf("Failed to insert session for user 2: %v", err)
	}
	janeSessionID, _ := result.LastInsertId()

	err = assignProjectRole(db, int(janeSessionID), 1, 3, 1)
	if err == nil {
		t.Errorf("expected error for lack of exec permission but got none")
	}
}
