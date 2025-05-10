package users

import (
	"brickedup/backend/utils"
	"testing"

	_ "modernc.org/sqlite"
)

func TestLogin(t *testing.T) {
	db := utils.SetupTest(t)
	defer db.Close()

	// Test valid login
	session, err := Login(db, "john.doe@example.com", "hashed_password_1")
    if err != nil {
        t.Fatal(err)
    }

	if session.SessionID < 0 {
        t.Fatalf("Valid login failed: invalid sessionid")
    }

    // Test invalid password
    _, err = Login(db, "user1@example.com", "wrongpassword")
    if err == nil {
        t.Fatal("Invalid password failed: should not be logged in")
    }

	// Test non-existent user
	_, err = Login(db, "nouser@example.com", "testpassword")
    if err == nil {
        t.Fatal("Non-existent user failed: should not be logged in")
    }

	// Test unverified user
	_, err = Login(db, "unverified@example.com", "password3")
    if err == nil {
        t.Fatal("Unverified user failed: should not be logged in")
    }
}
