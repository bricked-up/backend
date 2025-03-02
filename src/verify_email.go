package backend

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// VerifyUser verifies the user's email using the provided verification code.
// It removes all expired verification codes from the database and sets the user's
// verifyid to NULL if the verification code is correct and has not expired.
func VerifyUser(verificationCode string) error {
	// Remove expired verification codes
	db, err := sql.Open("sqlite", "bricked-up_prod.db")
	if err != nil {
		return fmt.Errorf("failed to launch the database")
	}

	_, err = db.Exec("DELETE FROM verification_codes WHERE expires_at < ?", time.Now())
	if err != nil {
		return fmt.Errorf("failed to remove expired verification codes: %v", err)
	}

	// Verify the provided verification code
	var userID int
	err = db.QueryRow("SELECT user_id FROM verification_codes WHERE code = ? AND expires_at >= ?", verificationCode, time.Now()).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("invalid verification code")
		}
		return fmt.Errorf("failed to verify code: %v", err)
	}

	// Set user's verifyid to NULL
	_, err = db.Exec("UPDATE user SET verifyid = NULL WHERE id = ?", userID)
	if err != nil {
		return fmt.Errorf("failed to update user verification status: %v", err)
	}
	db.Close()
	return nil
}
