package backend

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

// VerifyUser verifies the user's email using the provided verification code.
// It removes all expired verification codes from the database and sets the user's
// verifyid to NULL if the verification code is correct and has not expired.
func VerifyUser(verificationCode string, db *sql.DB) error {
	var err error
	_, err = db.Exec("DELETE FROM verification_codes WHERE expires_at < ?", time.Now())
	if err != nil {
		return err
	}

	// Verify the provided verification code
	var userID int
	err = db.QueryRow("SELECT user_id FROM verification_codes WHERE code = ? AND expires_at >= ?", verificationCode, time.Now()).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err
	}

	// Set user's verifyid to NULL
	_, err = db.Exec("UPDATE user SET verifyid = NULL WHERE id = ?", userID)
	if err != nil {
		return err
	}
	db.Close()
	return nil
}
