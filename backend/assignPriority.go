package backend

import (
	"database/sql"
	"errors"
	"log"

	_ "modernc.org/sqlite"
)

// AssignIssuePriority assigns a priority to an issue within a project.
func AssignIssuePriority(db *sql.DB, sessionid int, projectIssueID int, priorityID int) error {
	// Validate user's session and retrieve their user ID
	var userID int
	err := db.QueryRow("SELECT userid FROM SESSION WHERE id = ?", sessionid).Scan(&userID)
	if err != nil {
		log.Printf("DEBUG: Session ID %d not found", sessionid)
		return errors.New("Invalid session ID")
	}

	// Get the project ID associated with the issue
	var projectID int
	err = db.QueryRow("SELECT projectid FROM PROJECT_ISSUES WHERE id = ?", projectIssueID).Scan(&projectID)
	if err != nil {
		log.Printf("DEBUG: Project issue ID %d not found", projectIssueID)
		return errors.New("Invalid projectIssueId")
	}

	// Verify the priority belongs to the same project
	var priorityProjectID int
	err = db.QueryRow("SELECT projectid FROM PRIORITY WHERE id = ?", priorityID).Scan(&priorityProjectID)
	if err != nil || priorityProjectID != projectID {
		log.Printf("DEBUG: Priority ID %d does not belong to project ID %d", priorityID, projectID)
		return errors.New("Invalid priority for the project")
	}

	// Check if the user has write permissions
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

	// Update the issue priority
	result, err := db.Exec("UPDATE ISSUE SET priorityid = ? WHERE id = (SELECT issueid FROM PROJECT_ISSUES WHERE id = ?)", priorityID, projectIssueID)
	if err != nil {
		log.Printf("DEBUG: Failed to update issue priority for issue ID %d", projectIssueID)
		return errors.New("Failed to update issue priority")
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
