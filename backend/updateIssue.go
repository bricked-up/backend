package backend

import (
	"brickedup/backend/utils"
	"database/sql"
	"errors"
	"time"
)

// Issue represents the issue table schema
type Issue struct {
	ID         int        `json:"id"`
	Title      string     `json:"title"`
	Desc       string     `json:"desc"`
	TagID      *int       `json:"tagid"`
	PriorityID *int       `json:"priorityid"`
	Created    time.Time  `json:"created"`
	Completed  *time.Time `json:"completed"`
	Cost       int        `json:"cost"`
}

// UpdateIssueDetails updates all non-foreign key fields of an issue if the user has write access.
func UpdateIssueDetails(db *sql.DB, sessionID int, issueID int, issue Issue) error {
	// Sanitize inputs
	title := utils.SanitizeText(issue.Title, utils.TEXT)
	desc := utils.SanitizeText(issue.Desc, utils.TEXT)

	// Get user ID from session
	var userID int
	err := db.QueryRow("SELECT userid FROM SESSION WHERE id = ?", sessionID).Scan(&userID)
	if err != nil {
		return errors.New("invalid session ID")
	}

	// Check if issue exists and get project ID
	var projectID int
	err = db.QueryRow(`
		SELECT p.projectid
		FROM PROJECT_ISSUES p
		JOIN ISSUE i ON p.issueid = i.id
		WHERE i.id = ?
	`, issueID).Scan(&projectID)
	if err != nil {
		return errors.New("issue not found")
	}

	// Check if user has write access in that project
	var hasWrite bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM PROJECT_MEMBER_ROLE pmr
			JOIN PROJECT_MEMBER pm ON pmr.memberid = pm.id
			JOIN PROJECT_ROLE pr ON pmr.roleid = pr.id
			WHERE pm.userid = ? AND pm.projectid = ? AND pr.can_write = 1
		)
	`, userID, projectID).Scan(&hasWrite)
	if err != nil || !hasWrite {
		return errors.New("user does not have write permissions")
	}

	// Update only non-foreign key fields
	result, err := db.Exec(`
		UPDATE ISSUE
		SET title = ?, desc = ?, created = ?, completed = ?, cost = ?
		WHERE id = ?
	`, title, desc, issue.Created, issue.Completed, issue.Cost, issueID)
	if err != nil {
		return errors.New("failed to update issue")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("could not retrieve update result")
	}
	if rowsAffected == 0 {
		return errors.New("no rows were updated")
	}

	return nil
}
