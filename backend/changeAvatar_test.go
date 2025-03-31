package backend

import (
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

func TestUpdateUserAvatar(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Use the existing PolarBear file in the current directory
	avatarPath := "C:/Users/pratu/OneDrive/Desktop/BrickedUp/backend/PolarBear.PNG"

	// Check that file actually exists
	if _, err := os.Stat(avatarPath); os.IsNotExist(err) {
		t.Fatalf("Avatar file does not exist at: %s", avatarPath)
	}

	// Load schema and data
	initSQL, _ := os.ReadFile("../sql/init.sql")
	populateSQL, _ := os.ReadFile("../sql/populate.sql")
	db.Exec(string(initSQL))
	db.Exec(string(populateSQL))

	// Insert valid session for user 1
	res, _ := db.Exec(`INSERT INTO SESSION (userid, expires) VALUES (?, datetime('now', '+1 day'))`, 1)
	sessionID, _ := res.LastInsertId()

	// SUCCESS CASE
	err = UpdateUserAvatar(db, int(sessionID), avatarPath)
	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}

	// ALREADY USED avatar
	_, _ = db.Exec(`UPDATE USER SET avatar = ? WHERE id = 2`, avatarPath)
	err = UpdateUserAvatar(db, int(sessionID), avatarPath)
	if err == nil || err.Error() != "avatar path is already in use" {
		t.Errorf("Expected 'avatar path is already in use', got: %v", err)
	}

	// INVALID SESSION
	err = UpdateUserAvatar(db, 9999, avatarPath)
	if err == nil || err.Error() != "invalid session" {
		t.Errorf("Expected invalid session error, got: %v", err)
	}

	// UNVERIFIED USER (userID=4 from populate.sql)
	res, _ = db.Exec(`INSERT INTO SESSION (userid, expires) VALUES (?, datetime('now', '+1 day'))`, 4)
	unverifiedSession, _ := res.LastInsertId()
	err = UpdateUserAvatar(db, int(unverifiedSession), avatarPath)
	if err == nil || err.Error() != "user is not verified" {
		t.Errorf("Expected unverified user error, got: %v", err)
	}
}
