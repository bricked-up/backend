package backend

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

const testDBFile = "test_database.db"

func setupTestDB() (*sql.DB, error) {
	_ = os.Remove(testDBFile) // Remove old database if exists
	db, err := sql.Open("sqlite", testDBFile)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
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

		CREATE TABLE VERIFY_USER (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code INTEGER UNIQUE NOT NULL,
			expires DATE NOT NULL
		);

		CREATE TABLE USER (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			name TEXT NOT NULL,
			avatar TEXT,
			verifyid INTEGER,
			FOREIGN KEY (verifyid) REFERENCES VERIFY_USER(id) ON DELETE SET NULL
		);

		CREATE TABLE SESSION (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			userid INTEGER NOT NULL,
			expires DATE NOT NULL,
			FOREIGN KEY (userid) REFERENCES USER(id) ON DELETE CASCADE
		);

		CREATE TABLE PROJECT (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			orgid INTEGER NOT NULL,
			name TEXT NOT NULL,
			budget INTEGER NOT NULL,
			charter TEXT NOT NULL,
			archived BOOLEAN NOT NULL,
			FOREIGN KEY (orgid) REFERENCES ORGANIZATION(id) ON DELETE CASCADE
		);

		CREATE TABLE PROJECT_ROLE (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			projectid INTEGER NOT NULL,
			name TEXT NOT NULL,
			can_read BOOLEAN NOT NULL,
			can_write BOOLEAN NOT NULL,
			can_exec BOOLEAN NOT NULL,
			FOREIGN KEY (projectid) REFERENCES PROJECT(id) ON DELETE CASCADE
		);

		CREATE TABLE PROJECT_MEMBER (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			userid INTEGER NOT NULL,
			projectid INTEGER NOT NULL,
			FOREIGN KEY (userid) REFERENCES USER(id) ON DELETE CASCADE,
			FOREIGN KEY (projectid) REFERENCES PROJECT(id) ON DELETE CASCADE
		);

		CREATE TABLE PROJECT_MEMBER_ROLE (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			memberid INTEGER NOT NULL,
			roleid INTEGER NOT NULL,
			FOREIGN KEY (memberid) REFERENCES PROJECT_MEMBER(id) ON DELETE CASCADE,
			FOREIGN KEY (roleid) REFERENCES PROJECT_ROLE(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return nil, err
	}

	// Insert test data
	_, err = db.Exec(`
		INSERT INTO ORGANIZATION (name) VALUES ('Test Org');
		INSERT INTO USER (email, password, name, avatar, verifyid) VALUES ('usera@example.com', 'password', 'User A', NULL, 1);
		INSERT INTO USER (email, password, name, avatar, verifyid) VALUES ('userb@example.com', 'password', 'User B', NULL, 1);
		INSERT INTO PROJECT (orgid, name, budget, charter, archived) VALUES (1, 'Test Project', 1000, 'Test Charter', 0);
		INSERT INTO PROJECT_ROLE (projectid, name, can_read, can_write, can_exec) VALUES (1, 'Admin', 1, 1, 1);
		INSERT INTO PROJECT_ROLE (projectid, name, can_read, can_write, can_exec) VALUES (1, 'Member', 1, 0, 0);
		INSERT INTO PROJECT_MEMBER (userid, projectid) VALUES (1, 1);
		INSERT INTO PROJECT_MEMBER (userid, projectid) VALUES (2, 1);
		INSERT INTO PROJECT_MEMBER_ROLE (memberid, roleid) VALUES (1, 1);
		INSERT INTO PROJECT_MEMBER_ROLE (memberid, roleid) VALUES (2, 2);
		INSERT INTO VERIFY_USER (code, expires) VALUES (1234, '2025-01-01');
		INSERT INTO SESSION (userid, expires) VALUES (1, '2025-01-01');
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestPromoteUserRole(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)
	defer db.Close()

	// Test Case 1: Valid promotion
	err = promoteUserRole(db, 1, 2, 1, 1) // Promote user B to Admin role in Test Project
	assert.NoError(t, err)

	// Test Case 2: User B already has the specified role
	err = promoteUserRole(db, 1, 2, 1, 1) // User B is already Admin in the project
	assert.Error(t, err)
	assert.Equal(t, "user B already has the specified role", err.Error())

	// Test Case 3: User A lacks exec permissions
	err = promoteUserRole(db, 2, 1, 1, 1) // User A is now User B with no exec permissions
	assert.Error(t, err)
	assert.Equal(t, "user A lacks exec permissions", err.Error())

	// Test Case 4: User B is not part of the project
	err = promoteUserRole(db, 1, 3, 1, 1) // User C is not in the project
	assert.Error(t, err)
	assert.Equal(t, "user B is not part of the project", err.Error())

	// Test Case 5: User B is not validated
	_, err = db.Exec(`UPDATE USER SET verifyid = NULL WHERE email = 'userb@example.com'`)
	if err != nil {
		t.Fatal(err)
	}

	err = promoteUserRole(db, 1, 2, 1, 1) // User B is not validated anymore
	assert.Error(t, err)
	assert.Equal(t, "user B is not validated", err.Error())

	// Test Case 6: User B does not exist
	err = promoteUserRole(db, 1, 999, 1, 1) // Non-existing User
	assert.Error(t, err)
	assert.Equal(t, "user B does not exist", err.Error())
}
