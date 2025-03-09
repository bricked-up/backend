package backend

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func RequestPasswordReset(token string, userEmail string, db *sql.DB) (string, error) {
	var email string
	err := db.QueryRow("SELECT email FROM reset WHERE reset_token = ? AND reset_token_expires = 1", token).Scan(&email)
	if err == sql.ErrNoRows {
		return "", err
	} else if err != nil {
		return "", err
	}
	db.Close()
	return email, nil
}
