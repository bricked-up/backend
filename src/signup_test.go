package backend

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

// TestRegisterUser verifies the registration process
func TestRegisterUser(t *testing.T) {
	// Create an in-memory SQLite database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Load database schema
	initSQL, err := os.ReadFile("../sql/init.sql")
	if err != nil {
		t.Fatalf("Failed to read init.sql: %v", err)
	}
	if _, err := db.Exec(string(initSQL)); err != nil {
		t.Fatalf("Failed to execute init.sql: %v", err)
	}

	// Debugging: Check if tables exist
	t.Log("Checking tables in database...")
	rows, _ := db.Query("SELECT name FROM sqlite_master WHERE type='table';")
	for rows.Next() {
		var name string
		rows.Scan(&name)
		t.Log("Found table:", name)
	}
	rows.Close()

	// Test user registration
	email := "test@example.com"
	password := "securepassword"

	err = registerUser(db, email, password)
	if err != nil {
		t.Errorf("registerUser failed: %v", err)
		return
	}

	// Verify user exists
	var userID int
	err = db.QueryRow("SELECT id FROM USER WHERE email = ?", email).Scan(&userID)
	if err != nil {
		t.Fatalf("User not found in database: %v", err)
	}

	// Verify code exists
	var code string
	var expire time.Time
	err = db.QueryRow("SELECT code, expires FROM VERIFY_USER WHERE id = ?", userID).Scan(&code, &expire)
	if err != nil {
		t.Fatalf("Verification record not found: %v", err)
	}

	// Check code validity
	if code == "" {
		t.Errorf("Generated code is empty")
	}
	if expire.Before(time.Now()) {
		t.Errorf("Verification code should not be expired")
	}
}

// TestGenerateVerificationCode checks the validity of generated codes
func TestGenerateVerificationCode(t *testing.T) {
	code1 := generateVerificationCode()
	code2 := generateVerificationCode()

	if code1 == code2 {
		t.Errorf("Generated codes should be unique, but got identical codes")
	}
}
