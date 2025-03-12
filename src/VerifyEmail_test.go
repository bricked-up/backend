package backend

import (
	"database/sql"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

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

	_, err = db.Exec(schema)
	if err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	return db
}

func TestVerifyUser_Success(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO user (id, verifyid) VALUES (1, 123)")
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	_, err = db.Exec("INSERT INTO verification_codes (code, user_id, expires_at) VALUES (?, ?, ?)",
		"valid-code", 1, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("failed to insert verification code: %v", err)
	}

	// Test verification
	err = VerifyUser("valid-code", db)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Verify the result
	var verifyID sql.NullInt64
	err = db.QueryRow("SELECT verifyid FROM user WHERE id = 1").Scan(&verifyID)
	if err != nil {
		t.Fatalf("failed to query user: %v", err)
	}

	if verifyID.Valid {
		t.Errorf("expected verifyid to be NULL, got %v", verifyID.Int64)
	}
}

func TestVerifyUser_InvalidOrExpiredCode(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO user (id, verifyid) VALUES (1, 123)")
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	_, err = db.Exec("INSERT INTO verification_codes (code, user_id, expires_at) VALUES (?, ?, ?)",
		"expired-code", 1, time.Now().Add(-time.Hour))
	if err != nil {
		t.Fatalf("failed to insert verification code: %v", err)
	}

	// Test verification with expired code
	err = VerifyUser("expired-code", db)
	if err == nil {
		t.Errorf("expected error for expired code, got nil")
	}

	// Test verification with invalid code
	err = VerifyUser("invalid-code", db)
	if err == nil {
		t.Errorf("expected error for invalid code, got nil")
	}
}
