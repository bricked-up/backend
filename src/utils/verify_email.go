package utils

import (
	"database/sql"
	"fmt"
	"time"
)

// verifyUser verifies the user's email using the provided verification code.
// It removes all expired verification codes from the database and sets the user's
// verifyid to NULL if the verification code is correct and has not expired.
func verifyUser(verificationCode string) error {
	// Remove expired verification codes

	_, err := db.Exec("DELETE FROM verification_codes WHERE expires_at < ?", time.Now())
	if err != nil {
		return fmt.Errorf("failed to remove expired verification codes: %v", err)
	}

	// Verify the provided verification code
	var userID int
	err = db.QueryRow("SELECT user_id FROM verification_codes WHERE code = ? AND expires_at >= ?", verificationCode, time.Now()).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("invalid or expired verification code")
		}
		return fmt.Errorf("failed to verify code: %v", err)
	}

	// Set user's verifyid to NULL
	_, err = db.Exec("UPDATE user SET verifyid = NULL WHERE id = ?", userID)
	if err != nil {
		return fmt.Errorf("failed to update user verification status: %v", err)
	}

	return nil
}
