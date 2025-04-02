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

	// Test cases
	testCases := []struct {
		name          string
		sessionID     int
		issueID       int
		expectedError error
	}{
		{
			name:          "Successful Issue Closure",
			sessionID:     1, // This should be a session ID from your populate.sql with write access
			issueID:       1, // This should be a valid issue ID from your populate.sql
			expectedError: nil,
		},
		{
			name:          "Issue Does Not Exist",
			sessionID:     1,
			issueID:       999, // Non-existent issue ID
			expectedError: ErrIssueNotFound,
		},
		{
			name:          "Invalid Session",
			sessionID:     999, // Invalid session ID
			issueID:       1,
			expectedError: ErrInvalidSession,
		},
		{
			name:          "User Without Write Privileges",
			sessionID:     5, // This should be a session ID from your populate.sql without write access
			issueID:       1,
			expectedError: ErrInsufficientPrivileges,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Restore issue state before each test if needed (to handle test order independence)
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

				// Parse completed date and verify it's recent (within last minute)
				completedTime, err := time.Parse("2006-01-02 15:04:05", completed.String)
				if err != nil {
					t.Errorf("couldn't parse completed timestamp: %v", err)
					return
				}

				timeDiff := time.Since(completedTime)
				if timeDiff > time.Minute {
					t.Errorf("completed timestamp is not recent: %v", completedTime)
				}
			}
		})
	}
}
