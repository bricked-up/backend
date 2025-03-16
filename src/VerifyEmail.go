package backend

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

// VerifyUser verifies the user's email using the provided verification code.
// It removes expired codes and sets the user's `verified` field to true.
func VerifyUser(verificationCode int, db *sql.DB) error {
	// Remove expired verification codes
	_, err := db.Exec("DELETE FROM VERIFY_USER WHERE expires < ?", time.Now())
	if err != nil {
		return err
	}

	// Check if the verification code exists and is valid
	var userID int
	err = db.QueryRow(`
		SELECT u.id FROM USER u 
		INNER JOIN VERIFY_USER vu ON u.verifyid = vu.id 
		WHERE vu.code = ? AND vu.expires >= ?`, verificationCode, time.Now()).Scan(&userID)

	if err != nil {
		return err
	}

	// Mark the user as verified and remove the verifyid
	_, err = db.Exec("UPDATE USER SET verifyid = NULL, verified = 1 WHERE id = ?", userID)
	if err != nil {
		return err
	}

	return nil
}
