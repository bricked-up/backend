package backend

import (
	"database/sql"
	"testing"
	"time"

	_ "modernc.org/sqlite" // Import SQLite driver for in-memory database
)

// TestRegisterUser verifies the registration process
func TestRegisterUser(t *testing.T) {
	// Create an in-memory SQLite database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create necessary tables
	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE,
			password TEXT
		);
		CREATE TABLE verify_users (
			id INTEGER PRIMARY KEY,
			code TEXT,
			expire DATETIME
		);
	`)
	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	// Test user registration
	email := "test@example.com"
	password := "securepassword"

	err = registerUser(db, email, password)
	if err != nil {
		t.Errorf("registerUser failed: %v", err)
	}

	// Verify user exists
	var userID int
	err = db.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&userID)
	if err != nil {
		t.Errorf("User not found in database: %v", err)
	}

	// Verify code exists
	var code string
	var expire time.Time
	err = db.QueryRow("SELECT code, expire FROM verify_users WHERE id = ?", userID).Scan(&code, &expire)
	if err != nil {
		t.Errorf("Verification record not found: %v", err)
	}

	// Check code validity
	if len(code) != 32 {
		t.Errorf("Generated code length incorrect, got %d, expected 32", len(code))
	}
	if expire.Before(time.Now()) {
		t.Errorf("Verification code should not be expired")
	}
}

// TestGenerateVerificationCode checks the validity of generated codes
func TestGenerateVerificationCode(t *testing.T) {
	code1 := generateVerificationCode()
	code2 := generateVerificationCode()

	if len(code1) != 32 {
		t.Errorf("Expected code length of 32, got %d", len(code1))
	}
	if code1 == code2 {
		t.Errorf("Generated codes should be unique, but got identical codes")
	}
}
