package backend

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func TestPromoteUserRole(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:") // Use in-memory database
	assert.NoError(t, err)
	defer db.Close()

	// Load database schema
	initSQL, err := os.ReadFile("../sql/init.sql")
	assert.NoError(t, err)
	_, err = db.Exec(string(initSQL))
	assert.NoError(t, err)

	// Load initial data
	populateSQL, err := os.ReadFile("../sql/populate.sql")
	assert.NoError(t, err)
	_, err = db.Exec(string(populateSQL))
	assert.NoError(t, err)

	// Debugging: Check if tables exist
	t.Log("Checking tables in database...")
	rows, _ := db.Query("SELECT name FROM sqlite_master WHERE type='table';")
	for rows.Next() {
		var name string
		rows.Scan(&name)
		t.Log("Found table:", name)
	}
	rows.Close()

	// Test Case 1: Valid promotion
	err = promoteUserRole(db, 1, 2, 1, 1)
	assert.NoError(t, err)

	// Test Case 2: User B already has the specified role
	err = promoteUserRole(db, 1, 2, 1, 1)
	assert.Error(t, err)
	assert.Equal(t, "user B already has the specified role", err.Error())

	// Test Case 3: User A lacks exec permissions
	err = promoteUserRole(db, 2, 1, 1, 1)
	assert.Error(t, err)
	assert.Equal(t, "user A lacks exec permissions", err.Error())

	// Test Case 4: User B is not part of the project
	err = promoteUserRole(db, 1, 3, 1, 1)
	assert.Error(t, err)
	assert.Equal(t, "user B is not part of the project", err.Error())

	// Test Case 5: User B is not validated
	_, err = db.Exec(`UPDATE USER SET verifyid = NULL WHERE email = 'userb@example.com'`)
	assert.NoError(t, err)

	err = promoteUserRole(db, 1, 2, 1, 1)
	assert.Error(t, err)
	assert.Equal(t, "user B is not validated", err.Error())

	// Test Case 6: User B does not exist
	err = promoteUserRole(db, 1, 999, 1, 1)
	assert.Error(t, err)
	assert.Equal(t, "user B does not exist", err.Error())
}
