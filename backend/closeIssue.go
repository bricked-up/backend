package backend

import (
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

// Error definitions
var (
	ErrIssueNotFound          = errors.New("issue does not exist")
	ErrInvalidSession         = errors.New("invalid or expired session")
	ErrSessionVerification    = errors.New("error verifying session")
	ErrInsufficientPrivileges = errors.New("user does not have write privileges for this project")
	ErrIssueUpdate            = errors.New("error closing issue")
	ErrPrivilegeCheck         = errors.New("error checking user privileges")
	ErrIssueRetrieval         = errors.New("error retrieving issue")
)

// CloseIssue marks an issue as completed with the current timestamp
// It requires the user associated with the sessionID to have write privileges in the project containing the issue
func CloseIssue(db *sql.DB, sessionID int, issueID int) error {
	// Check if issue exists
	var projectID int
	err := db.QueryRow(`
		SELECT pi.projectid
		FROM PROJECT_ISSUES pi
		JOIN ISSUE i ON pi.issueid = i.id
		WHERE i.id = ?`, issueID).Scan(&projectID)

	if err != nil {
		if err == sql.ErrNoRows {
			return ErrIssueNotFound
		}
		return errors.Join(ErrIssueRetrieval, err)
	}

	// Get user ID from session
	var userID int
	err = db.QueryRow(`
		SELECT userid FROM SESSION 
		WHERE id = ? AND expires > datetime('now')`, sessionID).Scan(&userID)

	if err != nil {
		if err == sql.ErrNoRows {
			return ErrInvalidSession
		}
		return errors.Join(ErrSessionVerification, err)
	}

	// Check if user has write privileges for the project
	var hasWriteAccess bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 
			FROM PROJECT_MEMBER pm
			JOIN PROJECT_MEMBER_ROLE pmr ON pm.id = pmr.memberid
			JOIN PROJECT_ROLE pr ON pmr.roleid = pr.id
			WHERE pm.userid = ? AND pm.projectid = ? AND pr.can_write = 1
		)`, userID, projectID).Scan(&hasWriteAccess)

	if err != nil {
		return errors.Join(ErrPrivilegeCheck, err)
	}

	if !hasWriteAccess {
		return ErrInsufficientPrivileges
	}

	// Update issue to mark it as completed with current timestamp
	_, err = db.Exec(`
		UPDATE ISSUE 
		SET completed = datetime('now')
		WHERE id = ?`, issueID)

	if err != nil {
		return errors.Join(ErrIssueUpdate, err)
	}

	return nil
}
