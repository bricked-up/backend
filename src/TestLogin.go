package backend

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func TestLogin(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:") // Use in-memory database
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Load database schema
	initSQL, err := os.ReadFile("../sql/init.sql")
	if err != nil {
		t.Fatalf("Failed to read init.sql: %v", err)
	}
	if _, err := db.Exec(string(initSQL)); err != nil {
		t.Fatalf("Failed to execute init.sql: %v", err)
	}

	// Load initial data
	populateSQL, err := os.ReadFile("../sql/populate.sql")
	if err != nil {
		t.Fatalf("Failed to read populate.sql: %v", err)
	}
	if _, err := db.Exec(string(populateSQL)); err != nil {
		t.Fatalf("Failed to execute populate.sql: %v", err)
	}

	// Debugging: Check if tables exist
	t.Log("Checking tables in database...")
	rows, _ := db.Query("SELECT name FROM sqlite_master WHERE type='table';")
	for rows.Next() {
		var name string
		rows.Scan(&name)
		t.Log("Found table:", name)
	}
	rows.Close()

	// Test valid login
	expiresAt, err := login(db, "user1@example.com", "password1")
	assert.NoError(t, err)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour).UTC(), expiresAt.UTC(), time.Minute)

	// Test invalid password
	_, err = login(db, "user1@example.com", "wrongpassword")
	assert.Error(t, err)

	// Test non-existent user
	_, err = login(db, "nouser@example.com", "testpassword")
	assert.Error(t, err)

	// Test unverified user
	_, err = login(db, "unverified@example.com", "password3")
	assert.Error(t, err)

	// Test existing session reuse
	expTime := time.Now().Add(2 * time.Hour).UTC()
	_, err = db.Exec("INSERT INTO SESSION (userid, expires) VALUES (?, ?)", 1, expTime)
	assert.NoError(t, err)

	expiresAt, err = login(db, "user1@example.com", "password1")
	assert.NoError(t, err)
	assert.WithinDuration(t, expTime.UTC(), expiresAt.UTC(), time.Second)
}
