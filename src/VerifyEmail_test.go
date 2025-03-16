package backend

import (
	"database/sql"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

// setupTestDB initializes an in-memory SQLite database for testing purposes.
func setupVerifyEmailTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	schema := `
	PRAGMA foreign_keys = ON;


CREATE TABLE ORGANIZATION (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE VERIFY_USER (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code INTEGER UNIQUE NOT NULL,
    expires DATE NOT NULL
);

-- Tables that depend only on tables already created
CREATE TABLE USER (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    verifyid INTEGER,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    name TEXT NOT NULL,
    avatar TEXT,
    verified BOOLEAN NOT NULL DEFAULT 0,
    FOREIGN KEY (verifyid) REFERENCES VERIFY_USER(id) ON DELETE SET NULL
);

CREATE TABLE SESSION (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    userid INTEGER NOT NULL,
    expires TIMESTAMP NOT NULL,
    FOREIGN KEY (userid) REFERENCES USER(id) ON DELETE CASCADE
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

CREATE TABLE ORG_MEMBER (
    id INTEGER PRIMARY KEY,
    userid INTEGER NOT NULL,
    orgid INTEGER NOT NULL,
    FOREIGN KEY (userid) REFERENCES USER(id) ON DELETE CASCADE,
    FOREIGN KEY (orgid) REFERENCES ORGANIZATION(id) ON DELETE CASCADE
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

-- Tables dependent on PROJECT
CREATE TABLE PROJECT_ROLE (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    projectid INTEGER NOT NULL,
    name TEXT NOT NULL,
    can_read BOOLEAN NOT NULL,
    can_write BOOLEAN NOT NULL,
    can_exec BOOLEAN NOT NULL,
    FOREIGN KEY (projectid) REFERENCES PROJECT(id) ON DELETE CASCADE
);

CREATE TABLE TAG (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    projectid INTEGER NOT NULL,
    name TEXT NOT NULL,
    color TEXT NOT NULL,
    FOREIGN KEY (projectid) REFERENCES PROJECT(id) ON DELETE CASCADE
);

CREATE TABLE PRIORITY (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    projectid INTEGER NOT NULL,
    name TEXT NOT NULL,
    priority INTEGER NOT NULL,
    FOREIGN KEY (projectid) REFERENCES PROJECT(id) ON DELETE CASCADE
);

CREATE TABLE ISSUE (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    desc TEXT NOT NULL,
    tagid INTEGER,
    priorityid INTEGER,
    created TIMESTAMP NOT NULL,
    completed TIMESTAMP,
    cost INTEGER NOT NULL,
    FOREIGN KEY (tagid) REFERENCES TAG(id) ON DELETE SET NULL,
    FOREIGN KEY (priorityid) REFERENCES PRIORITY(id) ON DELETE SET NULL
);


CREATE TABLE DEPENDENCY (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    issueid INTEGER NOT NULL,
    dependency INTEGER NOT NULL,
    FOREIGN KEY (issueid) REFERENCES ISSUE(id) ON DELETE CASCADE,
    FOREIGN KEY (dependency) REFERENCES ISSUE(id) ON DELETE CASCADE
);

CREATE TABLE REMINDER (
    id INTEGER PRIMARY KEY,
    issueid INTEGER NOT NULL,
    userid INTEGER NOT NULL,
    FOREIGN KEY (issueid) REFERENCES ISSUE(id) ON DELETE CASCADE,
    FOREIGN KEY (userid) REFERENCES USER(id) ON DELETE CASCADE
);

CREATE TABLE ORG_MEMBER_ROLE (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberid INTEGER NOT NULL,
    roleid INTEGER NOT NULL,
    FOREIGN KEY (memberid) REFERENCES ORG_MEMBER(id) ON DELETE CASCADE,
    FOREIGN KEY (roleid) REFERENCES ORG_ROLE(id) ON DELETE CASCADE
);

CREATE TABLE ORG_PROJECTS (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    orgid INTEGER NOT NULL,
    projectid INTEGER NOT NULL,
    FOREIGN KEY (orgid) REFERENCES ORGANIZATION(id) ON DELETE CASCADE,
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

CREATE TABLE PROJECT_ISSUES (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    projectid INTEGER NOT NULL,
    issueid INTEGER NOT NULL,
    FOREIGN KEY (projectid) REFERENCES PROJECT(id) ON DELETE CASCADE,
    FOREIGN KEY (issueid) REFERENCES ISSUE(id) ON DELETE CASCADE
);

CREATE TABLE USER_ISSUES (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    userid INTEGER NOT NULL,
    issueid INTEGER NOT NULL,
    FOREIGN KEY (userid) REFERENCES USER(id) ON DELETE CASCADE,
    FOREIGN KEY (issueid) REFERENCES ISSUE(id) ON DELETE CASCADE,
    UNIQUE (userid, issueid)
);

CREATE TABLE FORGOT_PASSWORD (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    userid INTEGER NOT NULL,
    code INTEGER NOT NULL,
    expirationdate TIMESTAMP NOT NULL,
    FOREIGN KEY (userid) REFERENCES USER(id) ON DELETE CASCADE
);
CREATE TABLE VERIFICATION_CODES (
	id INTEGER PRIMARY KEY AUTOINCREMENT
);
`

	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	return db
}

func TestVerifyUser_Success(t *testing.T) {
	db := setupVerifyEmailTestDB(t)
	defer db.Close()

	// Insert verification code
	_, err := db.Exec("INSERT INTO VERIFY_USER (id, code, expires) VALUES (?, ?, ?)", 123, 111222, time.Now().Add(1*time.Hour))
	if err != nil {
		t.Fatalf("failed to insert verification code: %v", err)
	}

	// Insert test user
	_, err = db.Exec("INSERT INTO USER (id, verifyid, email, password, name) VALUES (?, ?, ?, ?, ?)",
		1, 123, "test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	// Verify the user using the code
	if err := VerifyUser(111222, db); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Check if verifyid is NULL and verified is true
	var verifyID sql.NullInt64
	var verified bool
	err = db.QueryRow("SELECT verifyid, verified FROM USER WHERE id = ?", 1).Scan(&verifyID, &verified)
	if err != nil {
		t.Fatalf("failed to query user: %v", err)
	}

	if verifyID.Valid {
		t.Errorf("expected verifyid to be NULL, got %v", verifyID.Int64)
	}

	if !verified {
		t.Errorf("expected user to be verified, got false")
	}
}
func TestVerifyUser_InvalidOrExpiredCode(t *testing.T) {
	db := setupVerifyEmailTestDB(t)
	defer db.Close()

	// Insert expired verification code
	_, err := db.Exec("INSERT INTO VERIFY_USER (id, code, expires) VALUES (?, ?, ?)", 123, 111222, time.Now().Add(-1*time.Hour))
	if err != nil {
		t.Fatalf("failed to insert expired verification code: %v", err)
	}

	// Insert test user
	_, err = db.Exec("INSERT INTO USER (id, verifyid, email, password, name) VALUES (?, ?, ?, ?, ?)",
		1, 123, "test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	// Try to verify with expired code
	if err := VerifyUser(111222, db); err == nil {
		t.Errorf("expected error for expired code, got nil")
	}

	// Try to verify with invalid code
	if err := VerifyUser(999999, db); err == nil {
		t.Errorf("expected error for invalid code, got nil")
	}
}
