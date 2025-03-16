package backend

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

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
