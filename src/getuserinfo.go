package backend

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	_ "modernc.org/sqlite"
)

// Session holds session details for the logged-in user.
type Session struct {
	UserID    int       `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

// getSessionDetails fetches session details for a user by email.
func getSessionDetails(db *sql.DB, email string) ([]byte, error) {
	var session Session
	var password string
	var verifyID int

	// Retrieve session data
	err := db.QueryRow(`SELECT id, password, verifyid FROM USER WHERE email = ?`, email).Scan(&session.UserID, &password, &verifyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Ensure the user is verified
	if verifyID == 0 {
		return nil, errors.New("user is not verified")
	}

	// Check for existing session
	err = db.QueryRow(`SELECT timestamp FROM SESSION WHERE userid = ? AND timestamp > ?`, session.UserID, time.Now()).Scan(&session.ExpiresAt)
	if err != nil {
		return nil, errors.New("no active session found")
	}

	// Convert session data to JSON
	jsonSession, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}

	return jsonSession, nil
}
