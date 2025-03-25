package backend

import (
	"database/sql"
	"errors"
	"log"
)

// UpdateIssue updates all non-foreign key fields in an existing issue
func UpdateIssue(db *sql.DB, sessionid int, issueid int, title string, desc string, created string, completed *string, cost int) error {
	// Validate the user's session and retrieve the corresponding user ID
	var userID int
	err := db.QueryRow("SELECT userid FROM SESSION WHERE id = ?", sessionid).Scan(&userID)
	if err != nil {
		log.Printf("DEBUG: Session ID %d not found in SESSION table", sessionid)
		return errors.New("Invalid session ID")
	}

	// Get the project ID associated with the issue
	var projectID int
	err = db.QueryRow("SELECT projectid FROM PROJECT_ISSUES WHERE issueid = ?", issueid).Scan(&projectID)
	if err != nil {
		log.Printf("DEBUG: Issue ID %d not found in PROJECT_ISSUES table", issueid)
		return errors.New("Invalid issue ID")
	}

	// Check if the user has write permissions for the project
	var hasWritePermission bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM PROJECT_MEMBER_ROLE AS pmr
			JOIN PROJECT_MEMBER AS pm ON pmr.memberid = pm.id
			JOIN PROJECT_ROLE AS pr ON pmr.roleid = pr.id
			WHERE pm.userid = ? AND pm.projectid = ? AND pr.can_write = 1
		)
	`, userID, projectID).Scan(&hasWritePermission)
	if err != nil || !hasWritePermission {
		log.Printf("DEBUG: User ID %d does not have write permissions for project ID %d", userID, projectID)
		return errors.New("User does not have write permissions")
	}

	// Update the issue's non-foreign key fields
	result, err := db.Exec(`
		UPDATE ISSUE 
		SET title = ?, desc = ?, created = ?, completed = ?, cost = ?
		WHERE id = ?
	`, title, desc, created, completed, cost, issueid)

	if err != nil {
		log.Printf("DEBUG: Failed to update issue ID %d", issueid)
		return errors.New("Failed to update issue")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("Failed to retrieve rows affected")
	}

	if rowsAffected == 0 {
		return errors.New("No rows updated")
	}

	return nil
}
