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
	_ = os.Remove(testDBFile)
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

		CREATE TABLE PROJECT (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			orgid INTEGER NOT NULL,
			name TEXT NOT NULL
		);

		CREATE TABLE PROJECT_ROLE (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			projectid INTEGER NOT NULL,
			name TEXT NOT NULL
		);

		CREATE TABLE PROJECT_MEMBER (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			userid INTEGER NOT NULL,
			projectid INTEGER NOT NULL
		);

		CREATE TABLE PROJECT_MEMBER_ROLE (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			memberid INTEGER NOT NULL,
			roleid INTEGER NOT NULL,
			FOREIGN KEY (memberid) REFERENCES PROJECT_MEMBER(id),
			FOREIGN KEY (roleid) REFERENCES PROJECT_ROLE(id)
		);

		CREATE TABLE VERIFY_USER (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code INTEGER UNIQUE NOT NULL
		);

		CREATE TABLE USER (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			verifyid INTEGER,
			FOREIGN KEY (verifyid) REFERENCES VERIFY_USER(id)
		);

		CREATE TABLE SESSION (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			userid INTEGER NOT NULL,
			FOREIGN KEY (userid) REFERENCES USER(id)
		);
	`)
	if err != nil {
		return nil, err
	}

	// Insert test data
	_, err = db.Exec(`
		INSERT INTO ORGANIZATION (name) VALUES ('Test Org');
		INSERT INTO USER (email, verifyid) VALUES ('usera@example.com', 1);
		INSERT INTO USER (email, verifyid) VALUES ('userb@example.com', 1);
		INSERT INTO PROJECT (orgid, name) VALUES (1, 'Test Project');
		INSERT INTO PROJECT_ROLE (projectid, name) VALUES (1, 'Admin');
		INSERT INTO PROJECT_ROLE (projectid, name) VALUES (1, 'Member');
		INSERT INTO PROJECT_MEMBER (userid, projectid) VALUES (1, 1);
		INSERT INTO PROJECT_MEMBER (userid, projectid) VALUES (2, 1);
		INSERT INTO PROJECT_MEMBER_ROLE (memberid, roleid) VALUES (1, 1);
		INSERT INTO PROJECT_MEMBER_ROLE (memberid, roleid) VALUES (2, 2);
		INSERT INTO VERIFY_USER (code) VALUES (1234);
		INSERT INTO SESSION (userid) VALUES (1);
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestRemoveUserRole(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)
	defer db.Close()

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
