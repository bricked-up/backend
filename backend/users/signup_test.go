package users

import (
	"brickedup/backend/utils"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

// TestRegisterUser verifies the registration process
func TestRegisterUser(t *testing.T) {
	db := utils.SetupTest(t)
	defer db.Close()

	// Test user registration
	email := "test@example.com"
	password := "securepassword"

	err := Signup(db, email, password)
	if err != nil {
		t.Errorf("Signup failed: %v", err)
		return
	}

	// Verify user exists
	var verifyUserID int
	err = db.QueryRow("SELECT verifyid FROM USER WHERE email = ?", email).Scan(&verifyUserID)
	if err != nil {
		t.Fatalf("User not found in database: %v", err)
	}

	// Verify code exists
	var code string
	var expire time.Time
	err = db.QueryRow("SELECT code, expires FROM VERIFY_USER WHERE id = ?", verifyUserID).Scan(&code, &expire)
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
