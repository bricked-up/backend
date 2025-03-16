package backend

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func TestRemoveUserRole(t *testing.T) {
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

	// Test valid role removal
	err = removeUserRole(db, "1", "2", 2, 1) // User A removes User B's role
	assert.NoError(t, err)

	// Test user B no longer has the role
	var roleExists bool
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM PROJECT_MEMBER_ROLE WHERE memberid = (SELECT id FROM PROJECT_MEMBER WHERE userid = 2 AND projectid = 1) AND roleid = 2)`).Scan(&roleExists)
	assert.NoError(t, err)
	assert.False(t, roleExists)

	// Test user B does not exist
	err = removeUserRole(db, "1", "999", 2, 1)
	assert.Error(t, err)
	assert.Equal(t, "user B does not exist", err.Error())

	// Test user A lacks exec permission
	err = removeUserRole(db, "2", "2", 2, 1) // User A lacks exec permission
	assert.Error(t, err)
	assert.Equal(t, "user A lacks exec permissions", err.Error())
}
