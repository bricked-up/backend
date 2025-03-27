package backend

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

// setupTestDB initializes an in-memory SQLite database for testing purposes.
func setupVerifyEmailTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}
	initSQL, err := os.ReadFile("../../init.sql")
	if err != nil {
		t.Fatalf("failed to open init.sql %v", err)
	}

	if _, err := db.Exec(string(initSQL)); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	populateSQL, err := os.ReadFile("../../populate.sql")
	if err != nil {
		t.Fatalf("failed to poppulate th")
	}
	if _, err := db.Exec(string(populateSQL)); err != nil {
		t.Fatalf("failed to populate database: %v", err)
	}
	return db

}

func TestVerifyUser_Success(t *testing.T) {
	db := setupVerifyEmailTestDB(t)
	defer db.Close()

	// Insert verification code
	_, err := db.Exec("INSERT INTO VERIFY_USER (id, code, expires) VALUES (?, ?, ?)", 123, 111222, time.Now().Add(1*time.Hour))
	if err != nil {
		t.Fatalf("failed to insert verification code: %v", err)
	}

	// Insert test user
	_, err = db.Exec("INSERT INTO USER (id, verifyid, email, password, name) VALUES (?, ?, ?, ?, ?)",
		1, 123, "test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	// Verify the user using the code
	if err := VerifyUser(111222, db); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Check if verifyid is NULL and verified is true
	var verifyID sql.NullInt64
	var verified bool
	err = db.QueryRow("SELECT verifyid, verified FROM USER WHERE id = ?", 1).Scan(&verifyID, &verified)
	if err != nil {
		t.Fatalf("failed to query user: %v", err)
	}

	if verifyID.Valid {
		t.Errorf("expected verifyid to be NULL, got %v", verifyID.Int64)
	}

	if !verified {
		t.Errorf("expected user to be verified, got false")
	}
}
func TestVerifyUser_InvalidOrExpiredCode(t *testing.T) {
	db := setupVerifyEmailTestDB(t)
	defer db.Close()

	// Insert expired verification code
	_, err := db.Exec("INSERT INTO VERIFY_USER (id, code, expires) VALUES (?, ?, ?)", 123, 111222, time.Now().Add(-1*time.Hour))
	if err != nil {
		t.Fatalf("failed to insert expired verification code: %v", err)
	}

	// Insert test user
	_, err = db.Exec("INSERT INTO USER (id, verifyid, email, password, name) VALUES (?, ?, ?, ?, ?)",
		1, 123, "test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	// Try to verify with expired code
	if err := VerifyUser(111222, db); err == nil {
		t.Errorf("expected error for expired code, got nil")
	}

	// Try to verify with invalid code
	if err := VerifyUser(999999, db); err == nil {
		t.Errorf("expected error for invalid code, got nil")
	}
}
