package backend

import (
	"database/sql"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func TestCloseIssue(t *testing.T) {
	// Open in-memory database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	defer db.Close()

	// Enable foreign keys
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
	}

	// Read and execute init.sql to create schema
	initSQL, err := os.ReadFile("../sql/init.sql")
	if err != nil {
		t.Fatalf("failed to read init.sql: %v", err)
	}

	// Split and execute init statements individually
	initStatements := strings.Split(string(initSQL), ";")
	for _, stmt := range initStatements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := db.Exec(stmt); err != nil {
			t.Fatalf("failed to execute init statement: %v\nStatement: %s", err, stmt)
		}
	}

	// Read and execute populate.sql to load test data
	populateSQL, err := os.ReadFile("../sql/populate.sql")
	if err != nil {
		t.Fatalf("failed to read populate.sql: %v", err)
	}

	// Split and execute populate statements individually
	populateStatements := strings.Split(string(populateSQL), ";")
	for _, stmt := range populateStatements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := db.Exec(stmt); err != nil {
			t.Fatalf("failed to execute populate statement: %v\nStatement: %s", err, stmt)
		}
	}

	// Create fresh sessions that won't be expired for our test
	setupTestSessions := `
		UPDATE SESSION SET expires = datetime('now', '+1 day') WHERE id IN (1, 2, 3, 4, 5);
	`
	_, err = db.Exec(setupTestSessions)
	if err != nil {
		t.Fatalf("failed to set up test sessions: %v", err)
	}

	// Find a valid admin session (a session with write privileges)
	var validAdminSessionID int
	err = db.QueryRow(`
		SELECT s.id FROM SESSION s
		JOIN USER u ON s.userid = u.id
		JOIN PROJECT_MEMBER pm ON u.id = pm.userid
		JOIN PROJECT_MEMBER_ROLE pmr ON pm.id = pmr.memberid
		JOIN PROJECT_ROLE pr ON pmr.roleid = pr.id
		WHERE s.expires > datetime('now') AND pr.can_write = 1
		LIMIT 1
	`).Scan(&validAdminSessionID)
	if err != nil {
		t.Log("Could not find a session with write access, setting up our own test data")
		setupOwnData := `
			INSERT INTO USER (id, email, password, name, verified) 
			VALUES (999, 'test@example.com', 'password', 'Test User', 1);
			
			INSERT INTO SESSION (id, userid, expires) 
			VALUES (999, 999, datetime('now', '+1 day'));
			
			INSERT INTO PROJECT_MEMBER (id, userid, projectid) VALUES (999, 999, 1);
			INSERT INTO PROJECT_MEMBER_ROLE (id, memberid, roleid) VALUES (999, 999, 1);
		`
		_, err = db.Exec(setupOwnData)
		if err != nil {
			t.Fatalf("failed to set up own test data: %v", err)
		}
		validAdminSessionID = 999
	}

	// Find a valid issue ID
	var validIssueID int
	err = db.QueryRow("SELECT id FROM ISSUE WHERE completed IS NULL LIMIT 1").Scan(&validIssueID)
	if err != nil {
		t.Fatalf("could not find a valid issue to test with: %v", err)
	}

	// Find a session without write privileges
	var nonAdminSessionID int
	err = db.QueryRow(`
		SELECT s.id FROM SESSION s
		JOIN USER u ON s.userid = u.id
		LEFT JOIN PROJECT_MEMBER pm ON u.id = pm.userid
		LEFT JOIN PROJECT_MEMBER_ROLE pmr ON pm.id = pmr.memberid
		LEFT JOIN PROJECT_ROLE pr ON pmr.roleid = pr.id
		WHERE s.expires > datetime('now') AND (pr.can_write IS NULL OR pr.can_write = 0)
		LIMIT 1
	`).Scan(&nonAdminSessionID)
	if err != nil {
		// Create a non-admin user without any write privileges
		setupNonAdminUser := `
			INSERT INTO USER (id, email, password, name, verified) 
			VALUES (998, 'nonadmin@example.com', 'password', 'Non-Admin User', 1);
			
			INSERT INTO SESSION (id, userid, expires) 
			VALUES (998, 998, datetime('now', '+1 day'));
		`
		_, err = db.Exec(setupNonAdminUser)
		if err != nil {
			t.Fatalf("failed to set up non-admin user: %v", err)
		}
		nonAdminSessionID = 998
	}

	// Test cases
	testCases := []struct {
		name          string
		sessionID     int
		issueID       int
		expectedError error
	}{
		{
			name:          "Successful Issue Closure",
			sessionID:     validAdminSessionID,
			issueID:       validIssueID,
			expectedError: nil,
		},
		{
			name:          "Issue Does Not Exist",
			sessionID:     validAdminSessionID,
			issueID:       999, // Non-existent issue ID
			expectedError: ErrIssueNotFound,
		},
		{
			name:          "Invalid Session",
			sessionID:     9999, // Invalid session ID
			issueID:       validIssueID,
			expectedError: ErrInvalidSession,
		},
		{
			name:          "User Without Write Privileges",
			sessionID:     nonAdminSessionID,
			issueID:       validIssueID,
			expectedError: ErrInsufficientPrivileges,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Restore issue state before each test if needed
			if tc.name == "Successful Issue Closure" {
				_, err := db.Exec("UPDATE ISSUE SET completed = NULL WHERE id = ?", tc.issueID)
				if err != nil {
					t.Fatalf("failed to reset issue state: %v", err)
				}
			}

			// Call the function
			err := CloseIssue(db, tc.sessionID, tc.issueID)

			// Check if the error matches what we expect
			if !errors.Is(err, tc.expectedError) && (err != nil || tc.expectedError != nil) {
				t.Errorf("expected error %v, got %v", tc.expectedError, err)
				return
			}

			// For successful case, verify that the issue was actually closed
			if tc.expectedError == nil {
				var completed sql.NullString
				err = db.QueryRow("SELECT completed FROM ISSUE WHERE id = ?", tc.issueID).Scan(&completed)
				if err != nil {
					t.Errorf("failed to retrieve updated issue: %v", err)
					return
				}
				if !completed.Valid {
					t.Errorf("issue not marked as completed")
					return
				}

				// Parse the timestamp formatted as "2006-01-02 15:04:05"
				completedTime, err := time.Parse(time.RFC3339, completed.String)
				if err != nil {
					// Try parsing using SQLite's default timestamp format
					completedTime, err = time.Parse("2006-01-02 15:04:05", completed.String)
					if err != nil {
						t.Errorf("couldn't parse completed timestamp: %v", err)
						return
					}
				}

				timeDiff := time.Since(completedTime)
				if timeDiff > time.Minute {
					t.Errorf("completed timestamp is not recent: %v", completed.String)
				}
			}
		})
	}
}
