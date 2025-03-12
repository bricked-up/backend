package backend

import (
	"database/sql"
	"fmt"
	"testing"

	_ "modernc.org/sqlite"
)

// Helper function to initialize the in-memory database schema
func initializeTestDB() (*sql.DB, error) {
	// Create an in-memory SQLite database using modernc driver
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}

	// Initialize only the necessary tables for testing
	schema := `
		PRAGMA foreign_keys = ON;

		CREATE TABLE ORGANIZATION (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL
		);

		CREATE TABLE ORG_ROLE (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			orgid INTEGER NOT NULL,
			name TEXT NOT NULL,
			can_read BOOLEAN NOT NULL,
			can_write BOOLEAN NOT NULL,
			can_exec BOOLEAN NOT NULL,
			FOREIGN KEY (orgid) REFERENCES ORGANIZATION(id) ON DELETE CASCADE
		);

		CREATE TABLE USER (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			name TEXT NOT NULL
		);

		CREATE TABLE SESSION (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			userid INTEGER NOT NULL,
			expires DATE NOT NULL,
			FOREIGN KEY (userid) REFERENCES USER(id) ON DELETE CASCADE
		);

		CREATE TABLE ORG_MEMBER (
			id INTEGER PRIMARY KEY,
			userid INTEGER NOT NULL,
			orgid INTEGER NOT NULL,
			FOREIGN KEY (userid) REFERENCES USER(id) ON DELETE CASCADE,
			FOREIGN KEY (orgid) REFERENCES ORGANIZATION(id) ON DELETE CASCADE
		);

		CREATE TABLE ORG_MEMBER_ROLE (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			memberid INTEGER NOT NULL,
			roleid INTEGER NOT NULL,
			FOREIGN KEY (memberid) REFERENCES ORG_MEMBER(id) ON DELETE CASCADE,
			FOREIGN KEY (roleid) REFERENCES ORG_ROLE(id) ON DELETE CASCADE
		);
	`

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Helper function to check if the user data was correctly inserted
func checkTestData(db *sql.DB) {
	var userCount int
	err := db.QueryRow("SELECT COUNT(*) FROM USER").Scan(&userCount)
	if err != nil {
		fmt.Println("Error querying USER table:", err)
	} else {
		fmt.Printf("User count: %d\n", userCount)
	}

	var orgCount int
	err = db.QueryRow("SELECT COUNT(*) FROM ORGANIZATION").Scan(&orgCount)
	if err != nil {
		fmt.Println("Error querying ORGANIZATION table:", err)
	} else {
		fmt.Printf("Organization count: %d\n", orgCount)
	}

	var roleCount int
	err = db.QueryRow("SELECT COUNT(*) FROM ORG_ROLE").Scan(&roleCount)
	if err != nil {
		fmt.Println("Error querying ORG_ROLE table:", err)
	} else {
		fmt.Printf("Role count: %d\n", roleCount)
	}
}

// Test for the case where everything works fine
func TestAssignOrgRole_NoError(t *testing.T) {
	// Initialize in-memory SQLite database
	db, err := initializeTestDB()
	if err != nil {
		t.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	// Insert test data into the database
	_, err = db.Exec(`
		INSERT INTO USER (email, password, name) VALUES ("userA@example.com", "password", "User A"), ("userB@example.com", "password", "User B");
		INSERT INTO ORGANIZATION (name) VALUES ("Org1");
		INSERT INTO ORG_MEMBER (userid, orgid) VALUES (1, 1), (2, 1);
		INSERT INTO ORG_ROLE (orgid, name, can_read, can_write, can_exec) VALUES (1, "Admin", true, true, true), (1, "Member", true, false, false);
		INSERT INTO ORG_MEMBER_ROLE (memberid, roleid) VALUES (1, 1), (2, 2);
		INSERT INTO SESSION (userid, expires) VALUES (1, "2025-12-31");
	`)
	if err != nil {
		t.Fatalf("Error inserting test data: %v", err)
	}

	// Check if data is correctly inserted
	checkTestData(db)

	// Call the function under test with updated arguments
	err = assignOrgRole(db, "userA@example.com", "Admin", 1, 1)

	// Assert no error occurred
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}
