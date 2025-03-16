package backend

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
			return nil, errors.New("User not found")
		}
		return nil, err
	}

	// Ensure the user is verified
	if verifyID == 0 {
		return nil, errors.New("User is not verified")
	}

	// Check for existing session
	err = db.QueryRow(`SELECT timestamp FROM SESSION WHERE userid = ? AND timestamp > ?`, session.UserID, time.Now()).Scan(&session.ExpiresAt)
	if err != nil {
		return nil, errors.New("No active session found")
	}

	// Convert session data to JSON
	jsonSession, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}

	return jsonSession, nil
}

// Test function for getSessionDetails
func TestGetSessionDetails(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	// Load schema
	schemaSQL, _ := os.ReadFile("../sql/init.sql")
	_, err = db.Exec(string(schemaSQL))
	assert.NoError(t, err)

	// Load sample data
	_, err = db.Exec(`
		INSERT INTO USER (id, email, password, verifyid, verified) VALUES (1, 'user1@example.com', 'password', 1, 1);
		INSERT INTO SESSION (userid, timestamp) VALUES (1, datetime('now', '+1 hour'));
	`)
	assert.NoError(t, err)

	// Test valid session retrieval
	jsonData, err := getSessionDetails(db, "user1@example.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Test unverified user
	_, err = db.Exec(`UPDATE USER SET verified = 0 WHERE id = 1`)
	_, err = getSessionDetails(db, "user1@example.com")
	assert.Error(t, err)
	assert.Equal(t, "User is not verified", err.Error())

	// Test user without session
	_, err = getSessionDetails(db, "nonexistent@example.com")
	assert.Error(t, err)
	assert.Equal(t, "User not found", err.Error())
}
