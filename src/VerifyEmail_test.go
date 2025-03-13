package backend

import (
	"database/sql"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

// setupTestDB initializes an in-memory SQLite database for testing purposes.
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	schema := `
	CREATE TABLE user (
		id INTEGER PRIMARY KEY,
		verifyid INTEGER
	);

	CREATE TABLE verification_codes (
		code TEXT,
		user_id INTEGER,
		expires_at DATETIME
	);`

	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	return db
}

func TestVerifyUser_Success(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test user
	if _, err := db.Exec("INSERT INTO user (id, verifyid) VALUES (?, ?)", 1, 123); err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	// Insert valid verification code
	if _, err := db.Exec("INSERT INTO verification_codes (code, user_id, expires_at) VALUES (?, ?, ?)",
		"valid-code", 1, time.Now().Add(1*time.Hour)); err != nil {
		t.Fatalf("failed to insert verification code: %v", err)
	}

	// Verify user with valid code
	if err := VerifyUser("valid-code", db); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Check if verifyid is set to NULL
	var verifyID sql.NullInt64
	if err := db.QueryRow("SELECT verifyid FROM user WHERE id = ?", 1).Scan(&verifyID); err != nil {
		t.Fatalf("failed to query user: %v", err)
	}

	if verifyID.Valid {
		t.Errorf("expected verifyid to be NULL, got %v", verifyID.Int64)
	}
}

func TestVerifyUser_InvalidOrExpiredCode(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test user
	if _, err := db.Exec("INSERT INTO user (id, verifyid) VALUES (?, ?)", 1, 123); err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	// Insert expired verification code
	if _, err := db.Exec("INSERT INTO verification_codes (code, user_id, expires_at) VALUES (?, ?, ?)",
		"expired-code", 1, time.Now().Add(-1*time.Hour)); err != nil {
		t.Fatalf("failed to insert expired verification code: %v", err)
	}

	// Verify user with expired code (expect error)
	if err := VerifyUser("expired-code", db); err == nil {
		t.Errorf("expected error for expired code, got nil")
	}

	// Verify user with invalid code (expect error)
	if err := VerifyUser("invalid-code", db); err == nil {
		t.Errorf("expected error for invalid code, got nil")
	}
}
