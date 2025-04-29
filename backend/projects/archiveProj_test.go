package projects

import (
	"brickedup/backend/utils"
	"testing"
	"time"
)

func TestArchiveProj(t *testing.T) {
	db := utils.SetupTest(t)
	defer db.Close()

	const projectid = 1

	// Insert session for user A (userID 1 - John, has exec permission in project 1)
	expiry := time.Now().Add(24 * time.Hour)
	res, err := db.Exec(`INSERT INTO SESSION (userid, expires) VALUES (?, ?)`, 1, expiry)
	if err != nil {
		t.Fatalf("Failed to insert session for user 1: %v", err)
	}
	sessionID, _ := res.LastInsertId()

	err = ArchiveProj(db, int(sessionID), 1)
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}

	var archived bool
	err = db.QueryRow(
		`SELECT archived FROM PROJECT WHERE id = ?`,
		projectid).Scan(&archived)

	if err != nil {
		t.Fatalf("Failed to query role: %v", err)
	}

	if !archived {
		t.Fatal("Project was not archived!")
	}

	// User does not exist
	err = ArchiveProj(db, -1, projectid)
	if err == nil {
		t.Fatal("ArchiveProj should fail when user does not exist!")
	}

	// Project does not exist
	err = ArchiveProj(db, int(sessionID), -1)
	if err == nil {
		t.Fatal("ArchiveProj should fail when project does not exist!")
	}

	// User has insufficient privileges (userID=2, Jane)
	res, err = db.Exec(`INSERT INTO SESSION (userid, expires) VALUES (?, ?)`, 2, expiry)
	if err != nil {
		t.Fatalf("Failed to insert session for user 2: %v", err)
	}
	noExecSessionID, _ := res.LastInsertId()

	err = ArchiveProj(db, int(noExecSessionID), projectid)
	if err == nil {
		t.Fatal("User with insufficient privileges should not be able to archive a project.")
	}
}
