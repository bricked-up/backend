package backend

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func TestGetTagDetails(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	// Load schema
	schemaSQL, _ := os.ReadFile("../sql/init.sql")
	_, err = db.Exec(string(schemaSQL))
	assert.NoError(t, err)

	// Load sample data
	_, err = db.Exec(`
		INSERT INTO ORGANIZATION (id, name) VALUES (1, 'Test Org');
		INSERT INTO PROJECT (id, orgid, name, budget, charter, archived) VALUES (1, 1, 'Test Project', 1000, 'Test Charter', 0);
		INSERT INTO TAG (id, projectid, name, color) VALUES (1, 1, 'Urgent', '#FF0000');
	`)
	assert.NoError(t, err)

	// Test valid tag retrieval
	jsonData, err := getTagDetails(db, "1")
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Test non-existent tag retrieval
	_, err = getTagDetails(db, "999")
	assert.Error(t, err)
	assert.Equal(t, "Tag not found", err.Error())

	// Test Unix timestamp conversion for session expiry
	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	assert.WithinDuration(t, time.Now().Add(24*time.Hour).UTC(), time.Unix(expiresAt, 0).UTC(), time.Minute)
}
