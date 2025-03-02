package backend

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

func RequestPasswordReset(token string, userEmail string) (string, error) {
	var email string
	var expiry time.Time

	db, err := sql.Open("sqlite", "bricked-up_prod.db")
	if err != nil {
		return "", err
	}

	err = db.QueryRow("SELECT email, reset_token_expires FROM users WHERE reset_token = ?", token).Scan(&email, &expiry)
	if err == sql.ErrNoRows {
		return "", err
	} else if err != nil {
		return "", err
	}

	if time.Now().After(expiry) {
		return "", err
	}
	db.Close()
	return email, nil
}
