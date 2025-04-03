package backend

import (
	"database/sql"
	"errors"
)

// Error definitions
var (
	ErrInvalidSession         = errors.New("invalid session")
	ErrSessionExpired         = errors.New("session expired")
	ErrIssueNotFound          = errors.New("issue not found")
	ErrInsufficientPrivileges = errors.New("user does not have write privileges for this project")
)

// CloseIssue marks an issue as completed by the current user
func CloseIssue(db *sql.DB, sessionID int, issueID int) error {
	// First check if the session is valid
	var userID int
	err := db.QueryRow(`
		SELECT userid FROM SESSION 
		WHERE id = ? AND expires > datetime('now')
	`, sessionID).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidSession
		}
		return err
	}

	// Next check if the issue exists
	var exists bool
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM ISSUE WHERE id = ?)`, issueID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrIssueNotFound
	}

	// Retrieve project ID for privilege check
	var projectID int
	err = db.QueryRow(`SELECT projectid FROM PROJECT_ISSUES WHERE issueid = ?`, issueID).Scan(&projectID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrIssueNotFound
		}
		return err
	}

	// Check if the user has write privileges using COUNT(*)
	var canWrite int
	err = db.QueryRow(`
		SELECT COUNT(*) FROM PROJECT_MEMBER pm
		JOIN PROJECT_MEMBER_ROLE pmr ON pm.id = pmr.memberid
		JOIN PROJECT_ROLE pr ON pmr.roleid = pr.id
		WHERE pm.userid = ? AND pm.projectid = ? AND pr.can_write = 1
	`, userID, projectID).Scan(&canWrite)
	if err != nil {
		return err
	}
	if canWrite == 0 {
		return ErrInsufficientPrivileges
	}

	// Use SQLite's strftime to format the current time as "2006-01-02 15:04:05"
	_, err = db.Exec(`
		UPDATE ISSUE 
		SET completed = strftime('%Y-%m-%d %H:%M:%S', 'now')
		WHERE id = ?
	`, issueID)

	return err
}
