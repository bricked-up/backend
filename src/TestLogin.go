package backend

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

const testDBFile = "test_database.db"

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

func setupTestDB() (*sql.DB, error) {
	_ = os.Remove(testDBFile) // Remove old database if exists
	db, err := sql.Open("sqlite", testDBFile)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		PRAGMA foreign_keys = ON;

		CREATE TABLE VERIFY_USER (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code INTEGER UNIQUE NOT NULL,
			expires TIMESTAMP NOT NULL
		);

		CREATE TABLE USER (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			name TEXT NOT NULL,
			avatar TEXT,
			verifyid INTEGER,
			FOREIGN KEY (verifyid) REFERENCES VERIFY_USER(id) ON DELETE SET NULL
		);

		CREATE TABLE SESSION (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			userid INTEGER NOT NULL,
			expires TIMESTAMP NOT NULL,
			FOREIGN KEY (userid) REFERENCES USER(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return nil, err
	}

	// Insert test users with hashed passwords
	passwords := map[string]string{
		"user1@example.com":      "password1",
		"user2@example.com":      "password2",
		"unverified@example.com": "password3",
	}

	for email, pass := range passwords {
		hashedPass, err := hashPassword(pass)
		if err != nil {
			return nil, err
		}
		_, err = db.Exec(`INSERT INTO USER (email, password, name, avatar, verifyid) VALUES (?, ?, 'Test User', NULL, ?)`, email, hashedPass, 1)
		if err != nil {
			return nil, err
		}
	}

	// Mark one user as unverified
	_, err = db.Exec(`UPDATE USER SET verifyid = NULL WHERE email = ?`, "unverified@example.com")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestLogin(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)
	defer db.Close()

	// Test valid login
	expiresAt, err := login(db, "user1@example.com", "password1")
	assert.NoError(t, err)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour).UTC(), expiresAt.UTC(), time.Minute)

	// Test invalid password
	_, err = login(db, "user1@example.com", "wrongpassword")
	assert.Error(t, err)

	// Test non-existent user
	_, err = login(db, "nouser@example.com", "testpassword")
	assert.Error(t, err)

	// Test unverified user
	_, err = login(db, "unverified@example.com", "password3")
	assert.Error(t, err)

	// Test existing session reuse
	expTime := time.Now().Add(2 * time.Hour).UTC()
	_, err = db.Exec("INSERT INTO SESSION (userid, expires) VALUES (?, ?)", 1, expTime)
	assert.NoError(t, err)

	expiresAt, err = login(db, "user1@example.com", "password1")
	assert.NoError(t, err)
	assert.WithinDuration(t, expTime.UTC(), expiresAt.UTC(), time.Second)
}
