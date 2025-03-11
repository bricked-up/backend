package backend

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// RequestPasswordReset checks if the reset token is valid and returns the email associated with the token.
// If the token is invalid, it returns an error.
func RequestPasswordReset(token string, userEmail string, db *sql.DB) (string, error) {
	var email string
	// Check if the reset token is valid and has not expited
	err := db.QueryRow("SELECT email FROM reset WHERE reset_token = ? AND reset_token_expires = 1", token).Scan(&email)
	if err == sql.ErrNoRows {
		return "", err
	} else if err != nil {
		return "", err
	}
	// Delete all expited reset tokens
	_, err = db.Exec("DELETE FROM reset WHERE reset_token_expires = ?", 0)
	if err != nil {
		return "", err
	}
	db.Close()
	return email, nil
}
