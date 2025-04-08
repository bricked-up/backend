package users

import (
	"brickedup/backend/utils"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"

	_ "modernc.org/sqlite"
)

// ValidateResetToken checks if the reset token is valid and returns the email associated with the token.
func ValidateResetToken(token string, newpassword string, db *sql.DB) (string, error) {
	var email string
	var expiresAt time.Time
	token = utils.SanitizeText(token, utils.TEXT)
	// Check if the reset token is valid and has not expired
	err := db.QueryRow("SELECT email, reset_token_expires FROM reset WHERE reset_token = ?", token).Scan(&expiresAt)
	if err == sql.ErrNoRows {
		return "", err
	} else if err != nil {
		return "", err
	}

	// Check if the token is expired
	if time.Now().After(expiresAt) {
		return "", sql.ErrNoRows
	}

	// Delete all expired reset tokens
	_, err = db.Exec("DELETE FROM reset WHERE reset_token_expires <= ?", time.Now())
	if err != nil {
		return "", err
	}

	// Update the user's password
	hash, err := bcrypt.GenerateFromPassword([]byte(newpassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	newpassword = string(hash)
	_, err = db.Exec("UPDATE USER SET password = ? WHERE email = ?", newpassword, email)
	if err != nil {
		return "", err
	}

	return email, nil
}
